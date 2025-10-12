@echo off
:: PathFinder Quick Scan Script
:: Usage: SCAN.bat <target-url>
::
:: Examples:
::   SCAN.bat honeypotlogs.com
::   SCAN.bat https://example.com
::   SCAN.bat example.com -concurrency 100
::

if "%1"=="" (
    echo.
    echo ===============================================================================
    echo  PathFinder - Quick Scan
    echo ===============================================================================
    echo.
    echo Usage: SCAN.bat ^<target-url^> [options]
    echo.
    echo Examples:
    echo   SCAN.bat honeypotlogs.com
    echo   SCAN.bat https://example.com
    echo   SCAN.bat example.com -concurrency 100 -timeout 5
    echo.
    echo Available Options:
    echo   -concurrency ^<n^>     Number of concurrent requests (default: 50)
    echo   -timeout ^<n^>         Timeout in seconds (default: 10)
    echo   -wordlist ^<file^>     Custom wordlist file (default: wordlist.txt)
    echo   -theme ^<name^>        Theme: matrix, rainbow, cyber, blood, skittles
    echo   -o ^<file^>            Output file
    echo   -of ^<format^>         Output format: text, json, csv
    echo.
    echo Press any key to exit...
    pause >nul
    exit /b 1
)

:: Extract target from first argument
set TARGET=%1
shift

:: Check if target has http:// or https://
echo %TARGET% | findstr /i "^http://" >nul
if %ERRORLEVEL% EQU 0 goto :run_scan

echo %TARGET% | findstr /i "^https://" >nul
if %ERRORLEVEL% EQU 0 goto :run_scan

:: Add https:// if no scheme provided
set TARGET=https://%TARGET%

:run_scan
echo.
echo ===============================================================================
echo  Starting PathFinder scan...
echo ===============================================================================
echo  Target: %TARGET%
echo  Wordlist: wordlist.txt
echo.
echo  Controls during scan:
echo    F1 - Cycle Theme
echo    F2 - Random Skittles (fun colors!)
echo    F3 - Globe Mode
echo    Q  - Quit
echo ===============================================================================
echo.

:: Run PathFinder with target and any additional arguments
pathfinder.exe -target %TARGET% %1 %2 %3 %4 %5 %6 %7 %8 %9

echo.
echo Scan complete. Check above for results.
echo.
pause
