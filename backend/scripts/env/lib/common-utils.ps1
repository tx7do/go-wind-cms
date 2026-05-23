<#
.SYNOPSIS
通用工具函数库
.DESCRIPTION
提供日志记录、错误处理等通用函数
.NOTES
编码: UTF-8 (NO BOM) | 兼容: PowerShell 5.1+
#>

function Log {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function ErrorLog {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function SuccessLog {
    param([string]$Message)
    Write-Host "[OK] $Message" -ForegroundColor Green
}

function InfoLog {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor White
}

function Initialize-ErrorHandling {
    <#
    .SYNOPSIS
    配置错误处理偏好
    .DESCRIPTION
    设置为 Continue，避免非致命错误中断脚本
    #>
    $script:ErrorActionPreference = 'Continue'
}

$global:CommonUtilsLoaded = $true
