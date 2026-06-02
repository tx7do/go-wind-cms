#!/usr/bin/env bash
################################################################################
##                    Go 环境安装函数库 (Unix)
##
## 安装 Go 运行时、配置环境变量
## 安全策略：已安装则跳过，不升级；下载失败不影响旧版本
## 使用方式：source "${SCRIPT_DIR}/lib/go-utils.sh"
##
## 编码: UTF-8 (LF) | 兼容: bash 4.0+
################################################################################

# 防止重复加载
if [ "${_GO_UTILS_LOADED:-0}" -eq 1 ]; then
  return 0 2>/dev/null || true
fi
_GO_UTILS_LOADED=1

# ============================================================================
# 架构检测
# ============================================================================

# 获取当前操作系统和架构（映射为 Go 官方下载路径格式）
get_go_os_arch() {
  local os=$(uname | tr '[:upper:]' '[:lower:]')
  local arch_raw=$(uname -m)
  local arch

  # 处理操作系统类型
  case "$os" in
      darwin|linux|freebsd|aix|dragonfly|illumos|openbsd|netbsd|plan9|solaris)
          # 支持的操作系统，继续处理架构
          ;;
      *)
          error "不支持的操作系统: $os"
          return 1
          ;;
  esac

  # 处理架构类型（映射为常见的架构名称）
  case "$arch_raw" in
      i386|i486|i586|i686|i786|386)
          arch="386"      # 32 位 x86 架构（老旧 Intel/AMD CPU）
        ;;
      x86_64|amd64)       # 64 位 x86 架构（Intel/AMD 主流 CPU）
          arch="amd64"
          ;;
      aarch64|arm64)
          arch="arm64"    # 64 位 ARM 架构（Apple Silicon、华为鲲鹏等）
          ;;
      armv7l|armv6l|arm)
          arch="arm"      # 32 位 ARM 架构（armv6/armv7，如树莓派早期型号）
          ;;
      s390x)
          arch="s390x"    # 64 位 IBM Z 架构（大型机）
          ;;
      ppc64)
          arch="ppc64"    # 64 位 PowerPC 架构（大端模式）
          ;;
      ppc64le)
          arch="ppc64le"  # 64 位 PowerPC 架构（小端模式）
          ;;
      loongarch64)
          arch="loong64"  # 64 位龙芯架构（国产龙芯 CPU）
          ;;
      riscv64)
          arch="riscv64"  # 64 位 RISC-V 架构（开源指令集，如各类 RISC-V 开发板）
          ;;
      mips64)
          arch="mips64"   # MIPS 64位（大端模式）
          ;;
      mips64el)
          arch="mips64el" # MIPS 64位（小端模式）
          ;;
      mips)
          arch="mips"     # MIPS 32位（大端模式）
          ;;
      mipsel|mipsle)
          arch="mipsle"   # MIPS 32位（小端模式）
          ;;
      *)
          error "不支持的架构: $arch_raw"
          return 1
          ;;
  esac

  echo "$os $arch"
}

# ============================================================================
# 版本查询
# ============================================================================

# 获取本地 Go 的版本（未安装则返回 "未安装"）
get_go_local_version() {
    if ! command -v go &> /dev/null; then
        echo "未安装"
        return
    fi
    go version | awk '{print $3}'
}

# 获取 Go 最新稳定版本号（通过官方 API）
# 成功: 输出版本号 (如 "go1.25.3")
# 失败: return 1 (不输出任何内容)
get_go_latest_version() {
    if ! command -v curl &> /dev/null; then
      warn "curl 未安装，跳过获取最新版本"
      return 1
    fi
    local result
    result=$(curl -sf https://go.dev/dl/?mode=json 2>/dev/null | jq -r '.[0].version' 2>/dev/null) || true
    if [ -n "$result" ] && [[ "$result" == go* ]]; then
      echo "$result"
    else
      return 1
    fi
}

# ============================================================================
# 安装函数
# ============================================================================

# 安装指定版本的 Go
# 安全策略:
#   - 先下载到临时目录，验证文件后再替换
#   - 替换前备份旧版本到 /usr/local/go.bak
#   - 下载或解压失败则中止并恢复，不影响旧版本
install_go_version() {
    local version=$1
    local os=$2
    local arch=$3
    local sudo_cmd=""
    local tmp_dir=""

    if [ "$EUID" -ne 0 ]; then
      sudo_cmd="sudo"
    fi

    # 使用临时目录下载
    tmp_dir=$(mktemp -d)
    local tarball="${tmp_dir}/${version}.${os}-${arch}.tar.gz"

    log "下载 Go $version for $os-$arch..."
    if ! wget -q -O "$tarball" "https://go.dev/dl/${version}.${os}-${arch}.tar.gz"; then
      rm -rf "$tmp_dir"
      error "下载 Go 失败，请检查网络连接。旧版本保持不变。"
      return 1
    fi

    # 验证下载文件非空
    if [ ! -s "$tarball" ]; then
      rm -rf "$tmp_dir"
      error "下载的 Go 安装包为空。旧版本保持不变。"
      return 1
    fi

    # 备份旧版本（如果存在）
    if [ -d /usr/local/go ]; then
      log "备份当前 Go 到 /usr/local/go.bak..."
      ${sudo_cmd} rm -rf /usr/local/go.bak
      ${sudo_cmd} mv /usr/local/go /usr/local/go.bak
    fi

    log "安装 Go $version..."
    if ! ${sudo_cmd} tar -C /usr/local -xzf "$tarball"; then
      # 解压失败，尝试恢复备份
      log "解压失败，尝试恢复备份..."
      if [ -d /usr/local/go.bak ]; then
        ${sudo_cmd} mv /usr/local/go.bak /usr/local/go
        log "已恢复旧版本 Go"
      fi
      rm -rf "$tmp_dir"
      error "解压 Go 安装包失败"
      return 1
    fi

    # 安装成功，清理备份和临时文件
    ${sudo_cmd} rm -rf /usr/local/go.bak
    rm -rf "$tmp_dir"

    log "更新 PATH..."
    export PATH=$PATH:/usr/local/go/bin
    local shell_rc="${TARGET_HOME:-$HOME}/.bashrc"
    if [[ "$SHELL" == *"zsh"* ]]; then
        shell_rc="${TARGET_HOME:-$HOME}/.zshrc"
    fi

    if ! grep -q "/usr/local/go/bin" "$shell_rc" 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> "$shell_rc"
        log "已将 Go 加入 $shell_rc，请执行 'source $shell_rc' 生效"
    else
        log "PATH 中已包含 Go 路径，跳过"
    fi

    success "Go $version 安装完成！"
    /usr/local/go/bin/go version
}

# ============================================================================
# 初始化函数（安全策略：已安装则跳过，不升级）
# ============================================================================

install_golang() {
  log "检查 Golang..."

  local local_version
  local_version=$(get_go_local_version)
  log "当前 Go 版本: $local_version"

  if [ "$local_version" != "未安装" ]; then
      # Go 已安装 - 安全策略：跳过，不升级
      success "Go 已安装: $(go version 2>/dev/null)"
      log "跳过 Go 安装/升级，不干扰现有环境"
      return 0
  fi

  # Go 未安装 - 执行安装
  local os_arch
  read -r OS ARCH <<< $(get_go_os_arch)
  log "操作系统: $OS, 架构: $ARCH"

  log "Go 未安装，开始安装最新版本..."
  local latest_version
  latest_version=$(get_go_latest_version) || latest_version="go1.25.3"
  log "将安装版本: $latest_version"

  install_go_version "$latest_version" "$OS" "$ARCH"

  log "Golang 安装完成！"
}
