$ErrorActionPreference = "Stop"

Write-Host "=> Installing Capivara..." -ForegroundColor Cyan

if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Go is not installed." -ForegroundColor Red
    exit 1
}

if (-not (Get-Command "git" -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Git is not installed." -ForegroundColor Red
    exit 1
}

$InstallDir = "$env:USERPROFILE\.local\bin"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notmatch [regex]::Escape($InstallDir)) {
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    $env:Path = "$env:Path;$InstallDir"
}

$TmpDir = Join-Path $env:TEMP "capivara_install_$([guid]::NewGuid().ToString().Substring(0,8))"
New-Item -ItemType Directory -Force -Path $TmpDir | Out-Null
Set-Location $TmpDir

git clone https://github.com/atendi9/capivara.git .
go build -o capivara.exe main.go
Move-Item -Path "capivara.exe" -Destination "$InstallDir\capivara.exe" -Force

Set-Location $env:USERPROFILE
Remove-Item -Path $TmpDir -Recurse -Force

Write-Host "=> Capivara installed at $InstallDir\capivara.exe!" -ForegroundColor Green
Write-Host "=> Installation complete. You can now run 'capivara' from your terminal." -ForegroundColor Green