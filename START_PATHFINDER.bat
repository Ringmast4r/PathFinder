@echo off
title PathFinder v3.0 ULTIMATE
color 0A
cls

echo.
echo ================================================================================
echo  PathFinder v3.0 ULTIMATE - Web Path Discovery Tool
echo ================================================================================
echo.
echo  Just enter a target URL and we'll scan it for you!
echo.
echo ================================================================================
echo.

:MENU
echo.
echo  [1] Quick Scan (Basic)
echo  [2] Fast Scan (High Speed)
echo  [3] Thorough Scan (Extensions + Filters)
echo  [4] Custom Command
echo  [5] View Help
echo  [6] Exit
echo.
set /p choice="Select option (1-6): "

if "%choice%"=="1" goto QUICK
if "%choice%"=="2" goto FAST
if "%choice%"=="3" goto THOROUGH
if "%choice%"=="4" goto CUSTOM
if "%choice%"=="5" goto HELP
if "%choice%"=="6" goto END
goto MENU

:QUICK
cls
echo.
echo ================================================================================
echo  QUICK SCAN MODE
echo ================================================================================
echo.
set /p target="Enter target (e.g., example.com or https://example.com): "
if "%target%"=="" (
    echo ERROR: No target provided!
    timeout /t 2 >nul
    goto MENU
)

REM Auto-add https:// if not present
echo %target% | findstr /i "http://" >nul
if %errorlevel% neq 0 (
    echo %target% | findstr /i "https://" >nul
    if %errorlevel% neq 0 (
        set target=https://%target%
    )
)

echo.
echo Starting quick scan of %target%...
echo.
pathfinder.exe -target %target% -wordlist wordlist.txt
echo.
echo ================================================================================
echo  SCAN COMPLETE!
echo ================================================================================
echo.
pause
goto MENU

:FAST
cls
echo.
echo ================================================================================
echo  FAST SCAN MODE (High Concurrency)
echo ================================================================================
echo.
set /p target="Enter target (e.g., example.com or https://example.com): "
if "%target%"=="" (
    echo ERROR: No target provided!
    timeout /t 2 >nul
    goto MENU
)

REM Auto-add https:// if not present
echo %target% | findstr /i "http://" >nul
if %errorlevel% neq 0 (
    echo %target% | findstr /i "https://" >nul
    if %errorlevel% neq 0 (
        set target=https://%target%
    )
)

echo.
echo Starting fast scan of %target% with 100 concurrent requests...
echo.
pathfinder.exe -target %target% -wordlist wordlist.txt -concurrency 100
echo.
echo ================================================================================
echo  SCAN COMPLETE!
echo ================================================================================
echo.
pause
goto MENU

:THOROUGH
cls
echo.
echo ================================================================================
echo  THOROUGH SCAN MODE (Extensions + Filtering)
echo ================================================================================
echo.
set /p target="Enter target (e.g., example.com or https://example.com): "
if "%target%"=="" (
    echo ERROR: No target provided!
    timeout /t 2 >nul
    goto MENU
)

REM Auto-add https:// if not present
echo %target% | findstr /i "http://" >nul
if %errorlevel% neq 0 (
    echo %target% | findstr /i "https://" >nul
    if %errorlevel% neq 0 (
        set target=https://%target%
    )
)

echo.
echo Starting thorough scan of %target%...
echo Looking for PHP, HTML, JS, TXT files...
echo Filtering for interesting status codes only...
echo.
pathfinder.exe -target %target% -wordlist wordlist.txt -x php,html,js,txt,asp,aspx -mc 200,301,302,401,403 -concurrency 75
echo.
echo ================================================================================
echo  SCAN COMPLETE!
echo ================================================================================
echo.
pause
goto MENU

:CUSTOM
cls
echo.
echo ================================================================================
echo  CUSTOM COMMAND MODE
echo ================================================================================
echo.
echo  Enter your full PathFinder command (without "pathfinder.exe")
echo  Example: -target https://example.com -wordlist wordlist.txt -x php,html
echo.
set /p customcmd="Command: "
if "%customcmd%"=="" (
    echo ERROR: No command provided!
    timeout /t 2 >nul
    goto MENU
)
echo.
echo Running: pathfinder.exe %customcmd%
echo.
pathfinder.exe %customcmd%
echo.
echo ================================================================================
echo  SCAN COMPLETE!
echo ================================================================================
echo.
pause
goto MENU

:HELP
cls
echo.
echo ================================================================================
echo  PathFinder v3.0 ULTIMATE - Help
echo ================================================================================
echo.
pathfinder.exe
echo.
echo ================================================================================
echo.
echo  Quick Examples:
echo.
echo  Basic:        -target https://example.com -wordlist wordlist.txt
echo  Extensions:   -target https://example.com -wordlist wordlist.txt -x php,html,js
echo  Filtering:    -target https://example.com -wordlist wordlist.txt -mc 200,301
echo  Fast:         -target https://example.com -wordlist wordlist.txt -concurrency 100
echo  Export JSON:  -target https://example.com -wordlist wordlist.txt -o results.json -of json
echo.
echo ================================================================================
echo.
pause
goto MENU

:END
cls
echo.
echo  Thanks for using PathFinder!
echo  Exiting...
echo.
timeout /t 2 >nul
exit
