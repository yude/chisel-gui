@echo off
title build
go build -ldflags="-s -w" -trimpath
echo Press any key to exit ...
timeout /t 6 > nul