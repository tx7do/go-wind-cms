$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# 1. 加载库
. "$ScriptDir\common-utils.ps1"

Initialize-ErrorHandling

# 2. 测试日志函数
Log "测试日志"
Warn "测试警告"
ErrorLog "测试错误"
SuccessLog "测试成功"
InfoLog "测试信息"

# 3. 验证标记
if ($global:CommonUtilsLoaded) { Write-Host "✓ 库加载成功" -ForegroundColor Green }

# 4. 测试 trap 是否生效（故意触发非致命错误）
Get-ChildItem "C:\NonExistentPath" -ErrorAction Stop