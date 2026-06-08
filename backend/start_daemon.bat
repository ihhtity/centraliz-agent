@echo off
setlocal enabledelayedexpansion

:: 守护进程脚本 - 自动重启服务
:: 当服务器意外关闭时自动重启

set "APP_NAME=Centraliz Backend"
set "APP_EXE=main.exe"
set "APP_DIR=%~dp0"
set "LOG_FILE=%APP_DIR%daemon.log"
set "RESTART_DELAY=3"
set "MAX_RESTARTS=10"

:: 切换到应用目录
cd /d %APP_DIR%

:: 初始化重启计数器
set "restart_count=0"

:: 检查是否已在运行
tasklist | find /i "%APP_EXE%" >nul
if not errorlevel 1 (
    echo [%date% %time%] 错误：%APP_EXE% 已经在运行中。
    echo [%date% %time%] 错误：%APP_EXE% 已经在运行中。 >> %LOG_FILE%
    pause
    exit /b 1
)

echo ============================================
echo %APP_NAME% 守护进程启动
echo 启动时间：%date% %time%
echo 应用目录：%APP_DIR%
echo 日志文件：%LOG_FILE%
echo ============================================

echo [%date% %time%] 守护进程已启动，监控 %APP_EXE% >> %LOG_FILE%

:START_LOOP
    :: 检查重启次数
    if %restart_count% geq %MAX_RESTARTS% (
        echo [%date% %time%] 警告：已达到最大重启次数 %MAX_RESTARTS%，停止重启 >> %LOG_FILE%
        echo ============================================
        echo 警告：已达到最大重启次数 %MAX_RESTARTS%，守护进程将退出
        echo ============================================
        pause
        exit /b 1
    )

    :: 编译Go代码（如果需要）
    if not exist "%APP_EXE%" (
        echo [%date% %time%] 编译应用程序...
        go build -o %APP_EXE% main.go
        if errorlevel 1 (
            echo [%date% %time%] 编译失败：%errorlevel% >> %LOG_FILE%
            echo 编译失败，按任意键退出...
            pause
            exit /b 1
        )
    )

    :: 启动应用程序
    echo [%date% %time%] 启动 %APP_NAME%...
    echo [%date% %time%] 启动 %APP_NAME%，重启次数：%restart_count% >> %LOG_FILE%
    
    start "" /WAIT "%APP_DIR%%APP_EXE%"
    
    :: 检查退出码
    set "exit_code=%errorlevel%"
    
    echo [%date% %time%] %APP_NAME% 已退出，退出码：%exit_code% >> %LOG_FILE%

    :: 判断是否需要重启
    :: 退出码 0：正常退出（如收到关闭信号）
    :: 退出码 1-255：异常退出，需要重启
    if %exit_code% equ 0 (
        echo [%date% %time%] %APP_NAME% 正常退出，守护进程结束 >> %LOG_FILE%
        echo ============================================
        echo %APP_NAME% 正常退出，守护进程结束
        echo ============================================
        exit /b 0
    ) else (
        set /a restart_count+=1
        echo [%date% %time%] %APP_NAME% 异常退出，等待 %RESTART_DELAY% 秒后重启（第 %restart_count% 次尝试） >> %LOG_FILE%
        echo ============================================
        echo %APP_NAME% 异常退出，退出码：%exit_code%
        echo 等待 %RESTART_DELAY% 秒后重启（第 %restart_count% 次尝试）
        echo ============================================
        
        :: 等待重启延迟
        timeout /t %RESTART_DELAY% /nobreak >nul
        
        goto START_LOOP
    )

endlocal