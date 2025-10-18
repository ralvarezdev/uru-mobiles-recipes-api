@echo off

REM Load environment variables from .env file if it exists
if exist .env (
    call :loadenv .env
)

REM Check if port is provided via environment variables
if "%PORT%"=="" (
    echo Error: No port specified.
    exit /b 1
)

REM Execute the Go binary on port specified
.\bin\server\server.exe -mode=prod -port=%PORT%
exit /b

:loadenv
for /f "usebackq delims=" %%A in (%1) do (
    echo %%A | findstr /b /v "#" >nul
    if not errorlevel 1 (
        set "%%A"
    )
)
exit /b