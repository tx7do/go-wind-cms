#!/usr/bin/env bash
################################################################################
##                    操作系统检测函数库 (Unix)
##
## 自动检测操作系统类型和包管理器
## 使用方式：source "${SCRIPT_DIR}/lib/os-utils.sh"
##
## 编码: UTF-8 (LF) | 兼容: bash 4.0+
################################################################################

# 防止重复加载
if [ "${_OS_UTILS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_OS_UTILS_LOADED=1

# ============================================================================
# 操作系统检测
# ============================================================================

# 检测操作系统类型和包管理器
# 返回值格式: "os_type|pkg_mgr|pkg_cmd|docker_setup"
detect_os_and_package_manager() {
  local os_type=""
  local pkg_mgr=""
  local pkg_cmd=""
  local docker_setup="linux"

  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    os_type="macos"
    pkg_mgr="brew"
    pkg_cmd="brew install"
    docker_setup="macos"
  elif [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "linux"* ]]; then
    # Linux 发行版检测
    if [ -f /etc/os-release ]; then
      . /etc/os-release
      case "$ID" in
        ubuntu|debian|linuxmint|pop)
          os_type="ubuntu"
          pkg_mgr="apt"
          pkg_cmd="${SUDO} apt-get install -y"
          ;;
        centos|rhel)
          os_type="centos"
          pkg_mgr="yum"
          pkg_cmd="${SUDO} yum install -y"
          ;;
        rocky|almalinux)
          os_type="rocky"
          pkg_mgr="dnf"
          pkg_cmd="${SUDO} dnf install -y"
          ;;
        fedora)
          os_type="fedora"
          pkg_mgr="dnf"
          pkg_cmd="${SUDO} dnf install -y"
          ;;
        *)
          log "检测到 Linux 发行版: $ID，根据包管理器自动选择..."
          if command -v apt-get >/dev/null 2>&1; then
            os_type="ubuntu"
            pkg_mgr="apt"
            pkg_cmd="${SUDO} apt-get install -y"
          elif command -v dnf >/dev/null 2>&1; then
            os_type="fedora"
            pkg_mgr="dnf"
            pkg_cmd="${SUDO} dnf install -y"
          elif command -v yum >/dev/null 2>&1; then
            os_type="centos"
            pkg_mgr="yum"
            pkg_cmd="${SUDO} yum install -y"
          else
            error "无法识别的 Linux 发行版，且未找到 apt-get/dnf/yum"
            exit 1
          fi
          ;;
      esac
    else
      # 旧系统可能没有 /etc/os-release，尝试其他方式
      if command -v apt-get >/dev/null 2>&1; then
        os_type="ubuntu"
        pkg_mgr="apt"
        pkg_cmd="${SUDO} apt-get install -y"
      elif command -v dnf >/dev/null 2>&1; then
        os_type="fedora"
        pkg_mgr="dnf"
        pkg_cmd="${SUDO} dnf install -y"
      elif command -v yum >/dev/null 2>&1; then
        os_type="centos"
        pkg_mgr="yum"
        pkg_cmd="${SUDO} yum install -y"
      else
        error "无法检测 Linux 发行版类型"
        exit 1
      fi
    fi
  else
    error "不支持的操作系统类型: $OSTYPE"
    exit 1
  fi

  echo "$os_type|$pkg_mgr|$pkg_cmd|$docker_setup"
}
