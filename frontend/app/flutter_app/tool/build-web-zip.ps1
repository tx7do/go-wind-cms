# Flutter Web build + zip packaging (deliverable parity with other frontends: build ends with a zip).
# Usage (from flutter_app dir): powershell -ExecutionPolicy Bypass -File tool/build-web-zip.ps1
# NOTE: keep this file ASCII-only. Windows PowerShell 5.1 reads .ps1 as ANSI unless a BOM is
# present, and UTF-8 Chinese comments get misread into mojibake that breaks parsing.
# Android apk is blocked until `flutter doctor --android-licenses` is accepted;
# web is the stable target for now.
$ErrorActionPreference = 'Stop'
Set-Location (Join-Path $PSScriptRoot '..')

flutter build web --release
if ($LASTEXITCODE -ne 0) { throw "flutter build web failed (exit $LASTEXITCODE)" }

# flutter build does not clean build/web; a stale archive left there from an
# earlier manual run would get double-packed into the deliverable. Remove it.
Get-ChildItem 'build/web' -Filter '*.zip' -ErrorAction SilentlyContinue | Remove-Item -Force

$zip = 'build/flutter-web.zip'
if (Test-Path $zip) { Remove-Item $zip }

$sevenZip = Get-Command 7z -ErrorAction SilentlyContinue
if ($sevenZip) {
    7z a -tzip -mx=6 '--' $zip 'build/web'
    if ($LASTEXITCODE -ne 0) { throw "7z packaging failed (exit $LASTEXITCODE)" }
} else {
    Compress-Archive -Path 'build/web' -DestinationPath $zip -Force
    if ($LASTEXITCODE -ne 0) { throw "Compress-Archive packaging failed (exit $LASTEXITCODE)" }
}

$sizeMb = [math]::Round((Get-Item $zip).Length / 1MB, 1)
Write-Host "[build-web-zip] created $zip ($sizeMb MB)"
