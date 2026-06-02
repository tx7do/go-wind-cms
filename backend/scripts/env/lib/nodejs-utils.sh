#!/usr/bin/env bash
################################################################################
##                    Node.js / PM2 安装函数库 (Unix)
##
## 安装 Node.js (NodeSource 22.x LTS) 和 PM2 进程管理器
## 安全策略：已安装则跳过，不覆盖已有安装
## 使用方式：source "${SCRIPT_DIR}/lib/nodejs-utils.sh"
##
## 编码: UTF-8 (LF) | 兼容: bash 4.0+
################################################################################

# 防止重复加载
if [ "${_NODEJS_UTILS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_NODEJS_UTILS_LOADED=1

# ============================================================================
# Node.js / PM2 安装
# ============================================================================

# 安装 Node.js 和 PM2
# 参数:
#   $1 - pkg_mgr: 包管理器 (brew|apt|yum|dnf)
#
# 安全策略:
#   - 如果 node 命令已存在 → 跳过 Node.js 安装，不覆盖
#   - 如果 pm2 命令已存在 → 跳过 PM2 安装，不覆盖
#   - 如果 pm2 服务已配置 → 跳过 startup 配置
install_nodejs_and_pm2() {
  local pkg_mgr=$1

  # ========== Node.js 安装 ==========
  if command -v node >/dev/null 2>&1; then
    log "检测到 Node.js 已安装: $(node -v 2>/dev/null || echo '版本未知')"
    log "跳过 Node.js 安装，不干扰现有环境"
  else
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
  fi

  # ========== PM2 安装 ==========
  if command -v pm2 >/dev/null 2>&1; then
    log "检测到 PM2 已安装: $(pm2 --version 2>/dev/null || echo '版本未知')"
    log "跳过 PM2 安装，不干扰现有环境"
  else
    log "全局安装 pm2..."
    ${SUDO} npm install -g pm2@latest || true
    pm2 --version || true
  fi

  # ========== PM2 开机启动配置 ==========
  # 检查 pm2 startup 是否已经配置过
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: 检查 launchd plist 是否存在
    local launchd_plist="$HOME/Library/LaunchAgents/pm2.${USER}.plist"
    if [ -f "$launchd_plist" ]; then
      log "PM2 开机启动已配置，跳过"
    else
      log "配置 pm2 开机启动（launchd）..."
      BREW_PREFIX="$(brew --prefix 2>/dev/null || echo '/opt/homebrew')"
      ${SUDO} env PATH="$PATH:${BREW_PREFIX}/bin" pm2 startup launchd -u "$USER" --hp "$HOME" || true
      pm2 save || true
    fi
  else
    # Linux: 检查 systemd 单元是否已存在
    if ${SUDO} systemctl status "pm2-${TARGET_USER}" >/dev/null 2>&1; then
      log "PM2 systemd 服务已配置，跳过"
    else
      log "为用户 ${TARGET_USER} 生成 pm2 systemd 单元..."
      ${SUDO} env HOME="${TARGET_HOME}" USER="${TARGET_USER}" PM2_HOME="${TARGET_HOME}/.pm2" PATH="/usr/bin:${PATH}" pm2 startup systemd -u "${TARGET_USER}" --hp "${TARGET_HOME}" || true
      ${SUDO} systemctl daemon-reload || true
      ${SUDO} systemctl enable --now "pm2-${TARGET_USER}" 2>/dev/null || true
    fi
  fi
}
