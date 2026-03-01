#!/usr/bin/env bash
####################################
## 安装Golang
####################################

set -euo pipefail

log() { echo "==> $*"; }
err_trap() { echo "Error: $*" >&2; exit 1; }

trap 'err_trap "安装 Golang 失败 (行 $LINENO)"' ERR

# 获取当前操作系统和架构
get_os_arch() {
  local os=$(uname | tr '[:upper:]' '[:lower:]')
  local arch_raw=$(uname -m)
  local arch

  # 处理操作系统类型
  case "$os" in
      darwin|linux|freebsd|aix|dragonfly|illumos|openbsd|netbsd|plan9|solaris)
          # 支持的操作系统，继续处理架构
          ;;
      *)
          err_trap "不支持的操作系统: $os"
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
          err_trap "不支持的架构: $arch_raw"
          ;;
  esac

  echo "$os $arch"
}

# 安装指定版本的Go
install_go() {
    local version=$1
    local os=$2
    local arch=$3
    local sudo_cmd=""

    if [ "$EUID" -ne 0 ]; then
      sudo_cmd="sudo"
    fi

    log "下载 Go $version for $os-$arch..."
    if ! wget -q https://go.dev/dl/$version.$os-$arch.tar.gz; then
      err_trap "下载 Go 失败，请检查网络连接"
    fi

    log "删除旧版本 Go..."
    ${sudo_cmd} rm -rf /usr/local/go

    log "安装 Go $version..."
    if ! ${sudo_cmd} tar -C /usr/local -xzf $version.$os-$arch.tar.gz; then
      err_trap "解压 Go 安装包失败"
    fi

    log "清理临时文件..."
    rm -f $version.$os-$arch.tar.gz

    log "更新 PATH..."
    export PATH=$PATH:/usr/local/go/bin
    if [[ "$SHELL" == *"zsh"* ]]; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
        source ~/.zshrc || true
    else
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        source ~/.bashrc || true
    fi

    log "Go $version 安装完成！"
    /usr/local/go/bin/go version
}

# 获取本地Go的版本
get_go_local_version() {
    if ! command -v go &> /dev/null; then
        echo "未安装"
        return
    fi
    go version | awk '{print $3}'
}

# 获取最新版本的Go
get_go_latest_version() {
    if ! command -v curl &> /dev/null; then
      log "curl 未安装，跳过获取最新版本"
      return 1
    fi
    curl -s https://go.dev/dl/?mode=json 2>/dev/null | jq -r '.[0].version' 2>/dev/null || echo "获取失败"
}

log "检查 Golang..."
GO_LOCAL_VERSION=$(get_go_local_version)
log "当前 Go 版本: $GO_LOCAL_VERSION"

read OS ARCH <<< $(get_os_arch)
log "操作系统: $OS, 架构: $ARCH"

if [ "$GO_LOCAL_VERSION" = "未安装" ]; then
    log "Go 未安装，开始安装最新版本..."
    GO_LATEST_VERSION=$(get_go_latest_version) || GO_LATEST_VERSION="go1.25.3"
    install_go $GO_LATEST_VERSION $OS $ARCH
elif GO_LATEST=$(get_go_latest_version); then
    if [ "$GO_LOCAL_VERSION" = "$GO_LATEST" ]; then
        log "Go 已是最新版本 ($GO_LOCAL_VERSION)"
    else
        log "发现新版本: $GO_LATEST (当前: $GO_LOCAL_VERSION)，开始升级..."
        install_go $GO_LATEST $OS $ARCH
    fi
else
    log "无法获取最新版本，跳过升级"
fi

log "Golang 检查/安装完成！"
