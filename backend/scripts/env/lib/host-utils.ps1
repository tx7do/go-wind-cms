<#
.SYNOPSIS
Hosts 文件管理工具函数库
.DESCRIPTION
提供 hosts 文件的增删改查功能
.NOTES
保存编码：UTF-8 with BOM | 需要管理员权限 | 兼容：PowerShell 5.1+
#>

# 导入通用工具库（如果尚未导入）
if (-not (Test-Path variable:global:CommonUtilsLoaded)) {
    $LibDir = Split-Path -Parent $MyInvocation.MyCommand.Path
    . "$LibDir\common-utils.ps1"
    Set-Variable -Name CommonUtilsLoaded -Value $true -Scope Global
}

function Edit-Hosts {
    <#
    .SYNOPSIS
    编辑 hosts 文件（添加或删除记录）
    .DESCRIPTION
    在系统 hosts 文件中添加或删除 IP 与域名的映射关系
    .PARAMETER IP
    IP 地址
    .PARAMETER Domain
    域名
    .PARAMETER Operate
    操作类型：Add（添加）或 Remove（删除）
    .EXAMPLE
    Edit-Hosts -IP "127.0.0.1" -Domain "postgres.local" -Operate "Add"
    .EXAMPLE
    Edit-Hosts -IP "127.0.0.1" -Domain "postgres.local" -Operate "Remove"
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]$IP,
        [Parameter(Mandatory=$true)]
        [string]$Domain,
        [ValidateSet("Add","Remove")]
        [string]$Operate = "Add"
    )

    $hostsFile = "$env:SystemRoot\System32\drivers\etc\hosts"

    # 校验管理员权限
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        ErrorLog "请以管理员身份运行脚本"
        return $false
    }

    try {
        $pattern = "^\s*$IP\s+$Domain\s*$"
        
        if ($Operate -eq "Add") {
            $content = Get-Content -Path $hostsFile -Raw -Encoding UTF8
            
            if ($content -match $pattern) {
                Log "记录已存在，无需重复添加: $IP $Domain"
                return $true
            }
            
            Add-Content -Path $hostsFile -Value "`n$IP $Domain" -Encoding UTF8
            SuccessLog "成功添加: $IP $Domain"
        }
        else {
            $lines = Get-Content -Path $hostsFile -Encoding UTF8
            $newLines = $lines | Where-Object { $_ -notmatch $pattern }
            
            if ($lines.Count -eq $newLines.Count) {
                Warn "记录不存在，无需删除: $IP $Domain"
                return $true
            }
            
            Set-Content -Path $hostsFile -Value $newLines -Encoding UTF8
            SuccessLog "成功移除: $IP $Domain"
        }

        # 刷新 DNS 缓存
        ipconfig /flushdns | Out-Null
        Log "DNS 缓存已刷新"
        
        return $true
    }
    catch {
        ErrorLog "操作失败: $($_.Exception.Message)"
        return $false
    }
}

function Initialize-Hosts {
    <#
    .SYNOPSIS
    批量初始化 hosts 记录
    .DESCRIPTION
    为多个服务批量添加 hosts 记录
    .PARAMETER Services
    服务名称数组
    .PARAMETER IP
    IP 地址（默认 127.0.0.1）
    .PARAMETER DomainSuffix
    域名后缀（默认 .local）
    .EXAMPLE
    Initialize-Hosts -Services @("postgres", "mysql", "redis") -IP "127.0.0.1"
    #>
    param(
        [Parameter(Mandatory=$true)]
        [string[]]$Services,
        [string]$IP = "127.0.0.1",
        [string]$DomainSuffix = ".local"
    )

    Log "========== 初始化 Hosts 记录 =========="
    
    $successCount = 0
    $failCount = 0
    
    foreach ($service in $Services) {
        $domain = "$service$DomainSuffix"
        $result = Edit-Hosts -IP $IP -Domain $domain -Operate "Add"
        
        if ($result) {
            $successCount++
        } else {
            $failCount++
        }
    }
    
    Log ""
    Log "完成: 成功 $successCount, 失败 $failCount"
    
    return ($failCount -eq 0)
}

