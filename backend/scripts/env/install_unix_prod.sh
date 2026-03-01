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

# ============================================================================
# 初始化和工具函数
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 日志函数
log() { echo "==> $*"; }
warn() { echo "⚠ WARNING: $*" >&2; }
err_trap() { echo "❌ 错误：第 $1 行发生错误" >&2; exit 1; }
trap 'err_trap $LINENO' ERR

# sudo 检查
SUDO=${SUDO:-}
if [ "$EUID" -ne 0 ]; then
  SUDO='sudo'
fi

# 目标用户和主目录
TARGET_USER=${SUDO_USER:-$(whoami)}
TARGET_HOME=$(eval echo "~${TARGET_USER}")

# ============================================================================
# 操作系统检测
# ============================================================================

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
            echo "❌ 错误：无法识别的 Linux 发行版，且未找到 apt-get/dnf/yum" >&2
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
        echo "❌ 错误：无法检测 Linux 发行版类型" >&2
        exit 1
      fi
    fi
  else
    echo "❌ 错误：不支持的操作系统类型: $OSTYPE" >&2
    exit 1
  fi

  echo "$os_type|$pkg_mgr|$pkg_cmd|$docker_setup"
}

# ============================================================================
# 主安装函数
# ============================================================================

install_basic_tools() {
  local os_type=$1
  local pkg_mgr=$2
  local pkg_cmd=$3

  log "安装基础工具..."

  case "$pkg_mgr" in
    brew)
      # macOS
      if ! command -v brew >/dev/null 2>&1; then
        log "Homebrew 未检测到，开始安装..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv || true)"
      fi
      log "更新 Homebrew..."
      brew update
      $pkg_cmd htop wget git jq unzip || true
      ;;
    apt)
      # Ubuntu/Debian
      export DEBIAN_FRONTEND=noninteractive
      log "更新软件包索引..."
      ${SUDO} apt-get update -y
      ${SUDO} apt-get upgrade -y
      $pkg_cmd htop wget unzip git jq ca-certificates curl gnupg lsb-release apt-transport-https software-properties-common make
      ;;
    yum)
      # CentOS
      export LC_ALL=C
      log "更新系统包索引..."
      ${SUDO} yum -y update || true
      ${SUDO} yum -y upgrade || true
      ${SUDO} yum -y install epel-release || true
      $pkg_cmd htop wget unzip git jq curl gnupg2 yum-utils
      ;;
    dnf)
      # Rocky/Fedora
      export LC_ALL=C
      log "检查并更新系统..."
      ${SUDO} dnf check-update || true
      ${SUDO} dnf -y upgrade --refresh

      # Rocky 特定：启用 CRB 仓库
      if grep -qi "rocky\|almalinux" /etc/os-release 2>/dev/null; then
        log "启用 CRB 仓库..."
        ${SUDO} dnf config-manager --set-enabled crb || true
        # 安装 EPEL
        ${SUDO} dnf -y install \
          https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm \
          https://dl.fedoraproject.org/pub/epel/epel-next-release-latest-9.noarch.rpm || true
      fi
      ${SUDO} dnf -y install epel-release htop wget unzip git jq curl gnupg2 dnf-plugins-core make
      ;;
  esac
}

install_nodejs_and_pm2() {
  local pkg_mgr=$1

  log "安装 Node.js (NodeSource 22.x LTS)..."

  case "$pkg_mgr" in
    brew)
      # macOS
      brew install node
      ;;
    apt)
      # Ubuntu/Debian
      curl -fsSL https://deb.nodesource.com/setup_22.x | ${SUDO} -E bash -
      ${SUDO} apt-get install -y nodejs
      ;;
    yum)
      # CentOS
      curl -fsSL https://rpm.nodesource.com/setup_22.x | ${SUDO} bash -
      ${SUDO} yum -y install nodejs
      ;;
    dnf)
      # Rocky/Fedora
      curl -fsSL https://rpm.nodesource.com/setup_22.x | ${SUDO} bash -
      ${SUDO} dnf -y install nodejs
      ;;
  esac

  log "Node 版本："
  node -v || true
  log "npm 版本："
  npm -v || true

  log "全局安装 pm2..."
  ${SUDO} npm install -g pm2@latest || true
  pm2 --version || true

  log "配置 pm2 开机启动..."
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: 使用 launchd
    BREW_PREFIX="$(brew --prefix)"
    log "配置 pm2 开机启动（launchd）..."
    ${SUDO} env PATH="$PATH:${BREW_PREFIX}/bin" pm2 startup launchd -u "$USER" --hp "$HOME" || true
    pm2 save || true
  else
    # Linux: 使用 systemd
    log "为用户 ${TARGET_USER} 生成 pm2 systemd 单元..."
    ${SUDO} env HOME="${TARGET_HOME}" USER="${TARGET_USER}" PM2_HOME="${TARGET_HOME}/.pm2" PATH="/usr/bin:${PATH}" pm2 startup systemd -u "${TARGET_USER}" --hp "${TARGET_HOME}" || true
    ${SUDO} systemctl daemon-reload || true
    ${SUDO} systemctl enable --now "pm2-${TARGET_USER}" 2>/dev/null || true
    # 尝试安装 pm2 bash 补全
    ${SUDO} env HOME="${TARGET_HOME}" USER="${TARGET_USER}" PATH="/usr/bin:${PATH}" bash -lc 'pm2 completion install >/dev/null 2>&1 || true'
  fi
}

install_docker() {
  local pkg_mgr=$1
  local docker_setup=$2

  if command -v docker >/dev/null 2>&1; then
    log "检测到 Docker 已安装，跳过安装步骤"
    return 0
  fi

  log "安装 Docker..."

  case "$docker_setup" in
    macos)
      # macOS
      if ! command -v docker >/dev/null 2>&1; then
        brew install --cask docker || true
      fi
      log "尝试启动 Docker Desktop..."
      open -a Docker || true

      # 等待 Docker 启动（最多等待 120 秒）
      log "等待 Docker 完全启动..."
      local n=0
      until docker info >/dev/null 2>&1; do
        n=$((n+1))
        if [ "$n" -ge 24 ]; then
          warn "等待 Docker 启动超时（120s）。请手动打开 Docker Desktop 并登录后重试。"
          break
        fi
        sleep 5
      done
      ;;
    linux)
      case "$pkg_mgr" in
        apt)
          # Ubuntu/Debian
          log "移除旧版 Docker..."
          for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
            ${SUDO} apt-get remove -y $pkg || true
          done

          log "配置 Docker 官方仓库..."
          ${SUDO} mkdir -p /etc/apt/keyrings
          curl -fsSL https://download.docker.com/linux/ubuntu/gpg | ${SUDO} gpg --dearmor -o /etc/apt/keyrings/docker.gpg
          ${SUDO} chmod a+r /etc/apt/keyrings/docker.gpg
          local arch=$(dpkg --print-architecture)
          local codename=$(. /etc/os-release && echo "$VERSION_CODENAME")
          echo "deb [arch=${arch} signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu ${codename} stable" | ${SUDO} tee /etc/apt/sources.list.d/docker.list > /dev/null

          ${SUDO} apt-get update -y
          ${SUDO} apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
          ;;
        yum)
          # CentOS
          log "移除旧版 Docker..."
          for pkg in docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine podman-docker containerd runc; do
            ${SUDO} yum -y remove $pkg >/dev/null 2>&1 || true
          done

          log "配置 Docker 官方仓库..."
          ${SUDO} yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || true
          ${SUDO} yum -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
          ;;
        dnf)
          # Rocky/Fedora
          log "移除旧版 Docker..."
          for pkg in docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine podman-docker containerd runc; do
            ${SUDO} dnf -y remove $pkg >/dev/null 2>&1 || true
          done

          log "配置 Docker 官方仓库..."
          ${SUDO} dnf -y install dnf-plugins-core
          ${SUDO} dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || true
          ${SUDO} dnf -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
          ;;
      esac

      log "启用并启动 Docker..."
      ${SUDO} systemctl enable --now docker || true

      log "将用户加入 docker 组..."
      ${SUDO} groupadd -f docker || true
      ${SUDO} usermod -aG docker "${TARGET_USER}" || true
      ;;
  esac

  log "Docker 状态检查完成。"
}

install_golang() {
  log "运行项目内的 Golang 安装脚本..."
  local golang_install_script="${SCRIPT_DIR}/install_golang.sh"

  if [ -f "${golang_install_script}" ]; then
    log "找到 Golang 安装脚本: ${golang_install_script}"
    if [ -x "${golang_install_script}" ]; then
      log "执行 Golang 安装脚本..."
      "${golang_install_script}"
    else
      log "脚本存在但不可执行，尝试用 bash 执行..."
      bash "${golang_install_script}"
    fi
  else
    warn "未找到 Golang 安装脚本: ${golang_install_script}，跳过"
  fi
}

cleanup() {
  local pkg_mgr=$1

  log "清理..."

  case "$pkg_mgr" in
    apt)
      ${SUDO} apt-get autoremove -y || true
      ${SUDO} apt-get autoclean -y || true
      ;;
    yum)
      ${SUDO} yum -y autoremove || true
      ${SUDO} yum clean all || true
      ;;
    dnf)
      ${SUDO} dnf -y autoremove || true
      ${SUDO} dnf clean all || true
      ;;
  esac
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

