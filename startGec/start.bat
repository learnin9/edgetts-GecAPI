@echo off
setlocal enabledelayedexpansion

:loop
start "" "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" "c:\gec\index.html"
timeout /t 5 /nobreak

:check_edge
tasklist /FI "IMAGENAME eq msedge.exe" | find /I /N "msedge.exe"
if "!errorlevel!"=="0" (
    timeout /t 25 /nobreak
    taskkill /im msedge.exe /t /f
    goto loop
) else (
    goto loop
)

:trap
echo 捕获到 Ctrl+C，正在关闭 Edge 浏览器...
taskkill /im msedge.exe /t /f
exit /b
