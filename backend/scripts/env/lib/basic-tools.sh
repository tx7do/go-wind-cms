#!/usr/bin/env bash
################################################################################
##                    基础工具安装函数库 (Unix)
##
## 安装基础工具：htop, wget, git, jq 等
## 使用方式：source "${SCRIPT_DIR}/lib/basic-tools.sh"
##
## 编码: UTF-8 (LF) | 兼容: bash 4.0+
################################################################################

# 防止重复加载
if [ "${_BASIC_TOOLS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_BASIC_TOOLS_LOADED=1

# ============================================================================
# 基础工具安装
# ============================================================================

# 安装基础工具
# 参数:
#   $1 - os_type: 操作系统类型 (macos|ubuntu|centos|rocky|fedora)
#   $2 - pkg_mgr: 包管理器 (brew|apt|yum|dnf)
#   $3 - pkg_cmd: 安装命令 (如 "sudo apt-get install -y")
# 逐个检查并安装缺失的工具（不升级已安装的包，避免破坏运行中的服务）
# 参数: $@ - 要检查/安装的工具名列表
_install_missing_tools() {
  local pkg_cmd=$1
  shift

  for tool in "$@"; do
    if command -v "$tool" >/dev/null 2>&1; then
      log "  [跳过] $tool 已安装"
    else
      log "  [安装] $tool 未找到，正在安装..."
      $pkg_cmd "$tool" || true
    fi
  done
}

install_basic_tools() {
  local os_type=$1
  local pkg_mgr=$2
  local pkg_cmd=$3

  log "检查基础工具..."

  case "$pkg_mgr" in
    brew)
      # macOS
      if ! command -v brew >/dev/null 2>&1; then
        log "Homebrew 未检测到，开始安装..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv || true)"
      else
        log "  [跳过] Homebrew 已安装"
      fi
      # brew install 对已安装的包会自动跳过，无需额外检测
      $pkg_cmd htop wget git jq unzip || true
      ;;
    apt)
      # Ubuntu/Debian - 仅更新索引和安装缺失的工具，不执行 upgrade
      export DEBIAN_FRONTEND=noninteractive
      log "更新软件包索引（不升级已安装的包）..."
      ${SUDO} apt-get update -y
      _install_missing_tools "$pkg_cmd" htop wget unzip git jq ca-certificates curl gnupg lsb-release apt-transport-https software-properties-common make
      ;;
    yum)
      # CentOS - 仅安装缺失的工具，不执行 upgrade
      export LC_ALL=C
      log "检查基础工具..."
      ${SUDO} yum -y install epel-release || true
      _install_missing_tools "$pkg_cmd" htop wget unzip git jq curl gnupg2 yum-utils
      ;;
    dnf)
      # Rocky/Fedora - 仅安装缺失的工具，不执行 upgrade
      export LC_ALL=C
      log "检查基础工具..."

      # Rocky 特定：启用 CRB 仓库
      if grep -qi "rocky\|almalinux" /etc/os-release 2>/dev/null; then
        log "启用 CRB 仓库..."
        ${SUDO} dnf config-manager --set-enabled crb || true
        ${SUDO} dnf -y install \
          https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm \
          https://dl.fedoraproject.org/pub/epel/epel-next-release-latest-9.noarch.rpm || true
      fi
      _install_missing_tools "$pkg_cmd" epel-release htop wget unzip git jq curl gnupg2 dnf-plugins-core make
      ;;
  esac
}
