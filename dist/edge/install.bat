@echo off
REM Liaison Edge 安装脚本 (Windows)
REM 此脚本会根据操作系统自动下载并安装对应的 Edge 安装包

setlocal enabledelayedexpansion

REM 默认配置
set "SERVER_HTTP_ADDR="
set "SERVER_EDGE_ADDR="
set "INSTALL_DIR=C:\Program Files\Liaison"
set "BIN_DIR=C:\Program Files\Liaison\bin"
set "CONFIG_DIR=C:\Program Files\Liaison\conf"
set "LOG_DIR=C:\Program Files\Liaison\logs"
set "BINARY_NAME=liaison-edge.exe"

REM 解析参数
set "ACCESS_KEY="
set "SECRET_KEY="

:parse_args
if "%~1"=="" goto end_parse
if "%~1"=="--help" goto show_help
if "%~1"=="-h" goto show_help

for /f "tokens=1,* delims==" %%a in ("%~1") do (
    set "arg_name=%%a"
    set "arg_value=%%b"
)

if "!arg_name!"=="--access-key" set "ACCESS_KEY=!arg_value!"
if "!arg_name!"=="--secret-key" set "SECRET_KEY=!arg_value!"
if "!arg_name!"=="--server-http-addr" set "SERVER_HTTP_ADDR=!arg_value!"
if "!arg_name!"=="--server-edge-addr" set "SERVER_EDGE_ADDR=!arg_value!"

shift
goto parse_args

:end_parse

REM 验证必需参数
if "%SERVER_HTTP_ADDR%"=="" (
    echo Error: --server-http-addr is required
    goto show_help
)
if "%SERVER_EDGE_ADDR%"=="" (
    echo Error: --server-edge-addr is required
    goto show_help
)
if "%ACCESS_KEY%"=="" (
    echo Error: --access-key is required
    goto show_help
)
if "%SECRET_KEY%"=="" (
    echo Error: --secret-key is required
    goto show_help
)

echo Starting Liaison Edge installation...
echo Detected OS/Arch: windows-amd64
echo Windows 安装路径: %INSTALL_DIR%

REM 创建临时目录
set "TMP_DIR=%TEMP%\liaison-edge-install"
if exist "%TMP_DIR%" rmdir /s /q "%TMP_DIR%"
mkdir "%TMP_DIR%"

REM 下载安装包
set "PACKAGE_NAME=liaison-edge-windows-amd64.tar.gz"
set "PACKAGE_URL=https://%SERVER_HTTP_ADDR%/packages/edge/%PACKAGE_NAME%"

echo Downloading installation package...
echo Package URL: %PACKAGE_URL%

REM 检查是否有 curl
where curl >nul 2>&1
if %errorlevel% equ 0 (
    curl -k -sSL -o "%TMP_DIR%\%PACKAGE_NAME%" "%PACKAGE_URL%"
    if !errorlevel! neq 0 (
        echo Error: Failed to download package
        exit /b 1
    )
) else (
    REM 尝试使用 PowerShell 下载（兼容 PowerShell 5.1）
    REM PowerShell 5.1 不支持 -SkipCertificateCheck，使用 ServerCertificateValidationCallback
    powershell -Command "& {[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; [System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}; Invoke-WebRequest -Uri '%PACKAGE_URL%' -OutFile '%TMP_DIR%\%PACKAGE_NAME%'}"
    if !errorlevel! neq 0 (
        echo Error: Failed to download package. Please install curl or use PowerShell.
        exit /b 1
    )
)

REM 解压安装包
echo Extracting package...
cd /d "%TMP_DIR%"

REM 检查是否有 tar (Windows 10 1803+ 自带)
where tar >nul 2>&1
if %errorlevel% equ 0 (
    tar -xzf "%PACKAGE_NAME%"
    if !errorlevel! neq 0 (
        echo Error: Failed to extract package with tar
        exit /b 1
    )
) else (
    REM tar 不可用，提示用户
    echo Error: tar command not found. Please use one of the following:
    echo   1. Update to Windows 10 1803 or later (includes tar)
    echo   2. Install Git for Windows (includes tar)
    echo   3. Install 7-Zip and manually extract the package
    echo   4. Use WSL or Git Bash to run install.sh instead
    exit /b 1
)

REM 检查二进制文件是否存在
if not exist "%TMP_DIR%\%BINARY_NAME%" (
    echo Error: Binary file %BINARY_NAME% not found in package
    exit /b 1
)

REM 创建安装目录
echo Installing...
if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

REM 复制二进制文件
copy "%TMP_DIR%\%BINARY_NAME%" "%BIN_DIR%\%BINARY_NAME%" >nul

REM 创建配置文件
echo Rendering configuration file...

REM 检查是否有模板文件
set "TEMPLATE_FILE="
if exist "%TMP_DIR%\liaison-edge.yaml.template" (
    set "TEMPLATE_FILE=%TMP_DIR%\liaison-edge.yaml.template"
)

if defined TEMPLATE_FILE (
    REM 从模板创建配置文件
    powershell -Command "(Get-Content '%TEMPLATE_FILE%') -replace '\$\{SERVER_ADDR\}', '%SERVER_EDGE_ADDR%' -replace '\$\{ACCESS_KEY\}', '%ACCESS_KEY%' -replace '\$\{SECRET_KEY\}', '%SECRET_KEY%' -replace '\$\{LOG_DIR\}', '%LOG_DIR:\=/%' | Set-Content '%CONFIG_DIR%\liaison-edge.yaml'"
) else (
    REM 创建默认配置文件
    (
        echo manager:
        echo   dial:
        echo     addrs:
        echo       - %SERVER_EDGE_ADDR%
        echo     network: tcp
        echo     tls:
        echo       enable: true
        echo       insecure_skip_verify: true
        echo   auth:
        echo     access_key: "%ACCESS_KEY%"
        echo     secret_key: "%SECRET_KEY%"
        echo log:
        echo   level: info
        echo   file: %LOG_DIR:\=/%/liaison-edge.log
        echo   maxsize: 100
        echo   maxrolls: 10
    ) > "%CONFIG_DIR%\liaison-edge.yaml"
)

echo Installation completed!
echo Edge binary: %BIN_DIR%\%BINARY_NAME%
echo Config file: %CONFIG_DIR%\liaison-edge.yaml
echo Edge will connect to: %SERVER_EDGE_ADDR%

REM 询问是否启动服务
echo.
echo Please choose how to run the service:
echo 1) Run as Windows Service (requires admin privileges)
echo 2) Run in background with nohup
echo 3) Skip, start manually later
echo.
set /p choice="Enter option [1-3] (default: 2): "
if "%choice%"=="" set "choice=2"

if "%choice%"=="1" (
    echo Setting up Windows Service...
    REM 这里可以添加 Windows 服务安装逻辑
    echo Windows Service setup not implemented yet. Please use option 2 or 3.
) else if "%choice%"=="2" (
    echo Starting Edge in background...
    start /b "" "%BIN_DIR%\%BINARY_NAME%" -c "%CONFIG_DIR%\liaison-edge.yaml" > "%LOG_DIR%\liaison-edge.log" 2>&1
    echo Edge started in background.
    echo Log file: %LOG_DIR%\liaison-edge.log
) else if "%choice%"=="3" (
    echo Skipping service setup.
    echo To start manually, run:
    echo   "%BIN_DIR%\%BINARY_NAME%" -c "%CONFIG_DIR%\liaison-edge.yaml"
)

REM 清理临时文件
rmdir /s /q "%TMP_DIR%" 2>nul

echo.
echo Installation completed successfully!

exit /b 0

:show_help
echo Usage: %~nx0 [OPTIONS]
echo.
echo Options:
echo   --access-key=KEY        Access key (required)
echo   --secret-key=KEY        Secret key (required)
echo   --server-http-addr=ADDR HTTP server address for downloading packages (required)
echo   --server-edge-addr=ADDR Edge server address for connection (required)
echo   --help, -h              Show this help message
echo.
echo Example:
echo   %~nx0 --access-key=xxx --secret-key=yyy --server-http-addr=example.com --server-edge-addr=example.com:30012
exit /b 1
