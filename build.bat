@echo off
set bd=%cd%\build
set configClient=config_client.json
set configServer=config_server.json

echo %bd%

rmdir /s /q %bd%
go build -o %bd%\klaatoo-explorer.exe %cd%\cmd\main.go
mklink /j %bd%\client %cd%\client
copy %cd%\%configClient% %bd%\%configClient%
copy %cd%\%configServer% %bd%\%configServer%