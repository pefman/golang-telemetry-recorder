@echo off
REM Windows Batch file to run F1 Telemetry Recorder
REM Double-click this file to start the application

echo ================================================
echo   Starting F1 Telemetry Recorder
echo ================================================
echo.

if not exist "f1-telemetry-recorder.exe" (
    echo ERROR: f1-telemetry-recorder.exe not found!
    echo.
    echo Please build the application first:
    echo   powershell -ExecutionPolicy Bypass -File build.ps1
    echo.
    pause
    exit /b 1
)

f1-telemetry-recorder.exe

pause
