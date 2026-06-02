#!/usr/bin/env bash
################################################################################
##                    Docker 安装函数库 (Unix)
##
## 安装 Docker：支持 macOS (Homebrew) 和 Linux (apt/yum/dnf)
## 安全策略：已安装则跳过，不卸载旧版，不破坏运行中的容器
## 使用方式：source "${SCRIPT_DIR}/lib/docker-utils.sh"
##
## 编码: UTF-8 (LF) | 兼容: bash 4.0+
################################################################################

# 防止重复加载
if [ "${_DOCKER_UTILS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_DOCKER_UTILS_LOADED=1

# ============================================================================
# Docker 安装
# ============================================================================

# 安装 Docker
# 参数:
#   $1 - pkg_mgr: 包管理器 (brew|apt|yum|dnf)
#   $2 - docker_setup: Docker 安装方式 (macos|linux)
#
# 安全策略:
#   - 如果 docker 命令已存在 → 跳过整个安装流程
#   - 不卸载任何旧版 Docker 包（避免杀掉运行中的容器）
#   - 仅在确认未安装时执行安装
install_docker() {
  local pkg_mgr=$1
  local docker_setup=$2

  # 第一道防线：检查 docker 命令是否可用
  if command -v docker >/dev/null 2>&1; then
    log "检测到 Docker 已安装: $(docker --version 2>/dev/null || echo '版本未知')"
    log "跳过 Docker 安装，不干扰现有环境"
    return 0
  fi

  # 第二道防线：Linux 下检查 docker 二进制是否存在于常见路径
  if [[ "$docker_setup" == "linux" ]]; then
    for docker_bin in /usr/bin/docker /usr/local/bin/docker /snap/bin/docker; do
      if [ -x "$docker_bin" ]; then
        log "检测到 Docker 二进制: $docker_bin"
        log "跳过 Docker 安装，不干扰现有环境"
        return 0
      fi
    done
  fi

  log "Docker 未安装，开始安装..."

  case "$docker_setup" in
    macos)
      # macOS
      brew install --cask docker || true
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
          # Ubuntu/Debian - 配置 Docker 官方仓库并安装
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
          log "配置 Docker 官方仓库..."
          ${SUDO} yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo || true
          ${SUDO} yum -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
          ;;
        dnf)
          # Rocky/Fedora
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

  log "Docker 安装完成。"
}
