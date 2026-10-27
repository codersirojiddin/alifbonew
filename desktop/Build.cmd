@echo off
cd /d "%~dp0"
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0build.ps1"
if errorlevel 1 (echo Build failed. & pause & exit /b 1)
echo Output: out\Harf.exe and out\Harf-Setup.exe
pause
