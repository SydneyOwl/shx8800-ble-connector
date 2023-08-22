@echo off

FOR /F %%v IN ('git describe --tags') DO SET BUILDVERSION=%%v
FOR /F %%c IN ('git rev-parse HEAD') DO SET COMMIT=%%c
SET BUILDTIME=%DATE%T%TIME%

md tmp

go mod tidy
set CC=x86_64-w64-mingw32-gcc

cd cmd\shell
REM 64
set GOARCH=amd64
go build -ldflags "-X 'github.com/sydneyowl/shx8800-ble-connector/config.VER=%BUILDVERSION%' -X 'github.com/sydneyowl/shx8800-ble-connector/config.COMMIT=%COMMIT%' -X 'github.com/sydneyowl/shx8800-ble-connector/config.BUILDTIME=%BUILDTIME%'"
move /Y shell.exe ..\..\tmp\shx8800-ble-connector_windows_amd64.exe 2>nul

cd ..\gui
go build -tags="gui" -ldflags "-X 'github.com/sydneyowl/shx8800-ble-connector/config.VER=%BUILDVERSION%' -X 'github.com/sydneyowl/shx8800-ble-connector/config.COMMIT=%COMMIT%' -X 'github.com/sydneyowl/shx8800-ble-connector/config.BUILDTIME=%BUILDTIME%' -H windowsgui"
move /Y gui.exe ..\..\tmp\shx8800-ble-connector-with-gui_windows_amd64.exe 2>nul