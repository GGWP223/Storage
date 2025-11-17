@echo off
cd /d %~dp0

start cmd /k "cd /d %~dp0cmd\auth && go run auth.go"
start cmd /k "cd /d %~dp0cmd\storage && go run storage.go"
start cmd /k "cd /d %~dp0cmd\web && go run web.go"