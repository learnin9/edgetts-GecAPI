@echo off
REM 设置 GOARCH 和 GOOS 环境变量
go env -w GOARCH=amd64
go env -w GOOS=windows

REM 读取当前版本号
set versionFile=version.txt

if not exist %versionFile% (
    echo v0.0.0 > %versionFile%
)

set /p version=<%versionFile%

REM 提取版本号中的主版本、次版本和修订版本
for /f "tokens=1-3 delims=." %%a in ("%version:v=%") do (
    set major=%%a
    set minor=%%b
    set patch=%%c
)

REM 叠加修订版本号
set /a patch=patch+1

REM 生成新的版本号
set newVersion=v%major%.%minor%.%patch%

REM 保存新的版本号到文件
echo %newVersion% > %versionFile%

REM 获取当前日期和时间并格式化为 YYYY-MM-DD HH:MM:SS
for /f "tokens=2 delims==" %%i in ('wmic os get localdatetime /value') do set buildDateTime=%%i
set buildDate=%buildDateTime:~0,4%-%buildDateTime:~4,2%-%buildDateTime:~6,2%
set buildTime=%buildDateTime:~8,2%:%buildDateTime:~10,2%:%buildDateTime:~12,2%
set buildDateTime=%buildDate% %buildTime%

REM 构建 Go 项目，并将版本号和构建时间注入到程序中
go build -ldflags "-s -w -X main.version=%newVersion% -X 'main.buildTime=%buildDateTime%'" .

REM 显示构建完成的消息
echo Build and compression completed successfully. Version: %newVersion%
