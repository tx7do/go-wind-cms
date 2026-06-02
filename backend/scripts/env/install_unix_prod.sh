#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

################################################################################
##                    统一 Unix/Linux 环境准备脚本
##
## 支持系统：macOS, Ubuntu/Debian, CentOS, Rocky/AlmaLinux, Fedora
## 自动检测操作系统并使用相应的包管理器
##
## 使用方式：
##   bash scripts/install_unix_prod.sh
##
################################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ============================================================================
# 加载函数库
# ============================================================================

_LIB_DIR="${SCRIPT_DIR}/lib"

# 按依赖顺序加载
source "${_LIB_DIR}/common-utils.sh"
source "${_LIB_DIR}/os-utils.sh"
source "${_LIB_DIR}/basic-tools.sh"
source "${_LIB_DIR}/docker-utils.sh"
source "${_LIB_DIR}/nodejs-utils.sh"
source "${_LIB_DIR}/go-utils.sh"
source "${_LIB_DIR}/host-utils.sh"

# 错误处理
trap 'err_trap $LINENO' ERR

# ============================================================================
# 清理函数（安全策略：仅清理缓存，不执行 autoremove）
# autoremove 可能删除服务运行依赖的库，在远程平台极其危险
# ============================================================================

cleanup() {
  local pkg_mgr=$1

  log "清理包管理器缓存..."

  case "$pkg_mgr" in
    apt)
      ${SUDO} apt-get autoclean -y || true
      ;;
    yum)
      ${SUDO} yum clean all || true
      ;;
    dnf)
      ${SUDO} dnf clean all || true
      ;;
    brew)
      brew cleanup --prune=30 2>/dev/null || true
      ;;
  esac
}

maybe_initialize_hosts() {
  local auto_init="${AUTO_INIT_HOSTS:-false}"

  if [[ ! "$auto_init" =~ ^(1|true|TRUE|yes|YES)$ ]]; then
    log "  [跳过] hosts 初始化未启用 (设置 AUTO_INIT_HOSTS=true 可开启)"
    return 0
  fi

  local hosts_ip="${HOSTS_IP:-127.0.0.1}"
  local hosts_domain_suffix="${HOSTS_DOMAIN_SUFFIX:-.local}"
  local hosts_services="${HOSTS_SERVICES:-postgres mysql redis}"

  read -r -a services <<< "$hosts_services"
  if [ "${#services[@]}" -eq 0 ]; then
    warn "HOSTS_SERVICES 为空，跳过 hosts 初始化"
    return 0
  fi

  log "初始化 hosts 记录..."
  initialize_hosts "$hosts_ip" "$hosts_domain_suffix" "${services[@]}"
}

# ============================================================================
# 主程序
# ============================================================================

main() {
  log "========================================"
  log "   Unix/Linux 环境准备脚本"
  log "========================================"
  log ""

  log "检测操作系统和包管理器..."
  local os_info=$(detect_os_and_package_manager)

  local os_type=$(echo "$os_info" | cut -d'|' -f1)
  local pkg_mgr=$(echo "$os_info" | cut -d'|' -f2)
  local pkg_cmd=$(echo "$os_info" | cut -d'|' -f3)
  local docker_setup=$(echo "$os_info" | cut -d'|' -f4)

  log "检测到系统: $os_type"
  log "包管理器: $pkg_mgr"
  log "目标用户: $TARGET_USER"
  log "用户主目录: $TARGET_HOME"
  log ""

  # 执行安装步骤
  install_basic_tools "$os_type" "$pkg_mgr" "$pkg_cmd"
  install_nodejs_and_pm2 "$pkg_mgr"
  install_docker "$pkg_mgr" "$docker_setup"
  install_golang
  maybe_initialize_hosts
  cleanup "$pkg_mgr"

  log ""
  log "========================================"
  log "   安装完成 ✓"
  log "========================================"
  log ""

  if [[ "$OSTYPE" == "darwin"* ]]; then
    log "建议重启终端以加载可能的环境变量更改。"
  else
    log "提示："
    log "  • 如果将用户加入 docker 组，需要重新登录以生效"
    log "  • pm2 的 systemd 单元已为用户 ${TARGET_USER} 启用"
    if grep -qi "docker" /etc/group 2>/dev/null && ! groups "${TARGET_USER}" | grep -q docker; then
      log "  • 请运行: newgrp docker 或重新登录以加入 docker 组"
    fi
  fi
  log ""
}

# 执行主程序
main
