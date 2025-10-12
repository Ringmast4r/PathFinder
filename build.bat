@echo off
echo ===============================================================================
echo  PathFinder - Building Go Edition
echo ===============================================================================
echo.

echo [1/3] Cleaning old builds...
if exist pathfinder.exe del pathfinder.exe
if exist pathfinder-windows.exe del pathfinder-windows.exe
if exist pathfinder-linux del pathfinder-linux
if exist pathfinder-mac del pathfinder-mac

echo [2/3] Building for Windows...
go build -ldflags="-s -w" -o pathfinder.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Build failed!
    exit /b 1
)
echo SUCCESS: pathfinder.exe created

echo [3/3] Building for Linux and Mac...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o pathfinder-linux main.go
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o pathfinder-mac main.go
set GOOS=windows
set GOARCH=amd64

echo.
echo ===============================================================================
echo  BUILD COMPLETE!
echo ===============================================================================
echo.
echo Created executables:
echo  - pathfinder.exe (Windows)
echo  - pathfinder-linux (Linux)
echo  - pathfinder-mac (macOS)
echo.
dir /b pathfinder*
echo.
echo Try it: pathfinder.exe -target https://example.com -wordlist wordlist.txt
echo ===============================================================================
