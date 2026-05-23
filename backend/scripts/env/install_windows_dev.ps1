<#
.SYNOPSIS
Windows 开发环境自动配置脚本（Scoop/Docker/Go）
#>

param(
    [switch]$SkipDocker,
    [switch]$AutoConfirm
)

#Set-StrictMode -Version Latest
$ErrorActionPreference = 'Continue'  # 改为Continue，避免非致命错误直接退出
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# 显示版本信息
Write-Host "PowerShell Version: $($PSVersionTable.PSVersion)" -ForegroundColor Cyan

# ========== 导入函数库（先加载，再配置） ==========
$libFiles = @(
    'common-utils.ps1',
    'scoop-utils.ps1', 
    'docker-utils.ps1',
    'go-utils.ps1',
    'host-utils.ps1'
)

foreach ($lib in $libFiles) {
    $libPath = Join-Path $ScriptDir "lib\$lib"
    if (Test-Path $libPath) {
        . $libPath
        Log "Loaded: $lib"
    } else {
        ErrorLog "Missing library: $libPath"
        exit 1
    }
}

Log "库加载测试"
Initialize-ErrorHandling

# 检测管理员权限
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$IsAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $IsAdmin) {
    Warn "Not running as administrator! Docker auto-start and hosts configuration will be skipped."
}

# ========== Hosts 配置（需要管理员权限）
if ($IsAdmin) {
    $services = @('postgres', 'mysql', 'redis', 'minio')
    Initialize-Hosts -Services $services -IP "127.0.0.1" -DomainSuffix ".local"
} else {
    Warn "Skipping hosts configuration (requires administrator privileges)"
}

# ========== Scoop 安装
Initialize-Scoop -Buckets @('main', 'extras') -Packages @('wget', 'unzip', 'git', 'jq', 'make', 'grep', 'gawk', 'sed', 'touch', 'mingw', 'nodejs', 'go')

# ========== Docker 安装
Initialize-Docker -SkipDocker $SkipDocker -IsAdmin $IsAdmin

# ========== Go 环境配置
Initialize-GoEnvironment -GoPath (Join-Path $env:USERPROFILE "go") -GoProxy "https://goproxy.io,direct" -SkipPlugins $false -SkipCliTools $false

# ========== 手动配置提示
$goPathValue = $env:GOPATH
if (-not $goPathValue) { $goPathValue = Join-Path $env:USERPROFILE "go" }
Log "Environment setup completed (current session only)!"
$tips = @"
==================== MANUAL CONFIG TIPS ====================
1. To make GOPATH permanent:
   Add these lines to your PowerShell Profile:
   `$env:GOPATH = "$goPathValue"
   `$env:PATH += ";$goPathValue\bin"

2. Check service status (admin):
   Get-Service com.docker.service
"@
Write-Host $tips -ForegroundColor Green
