echo off
@if "%1"=="" (
  echo usage: %0 source_filename[.go]
  exit /b

)
set filename=%~n1
set GOOS_TMP=%GOOS%
set GOARCH_TMP=%GOARCH%
set GOOS=windows
set GOARCH=amd64
go build -o %filename%.exe %filename%.go
set GOOS=%GOOS_TMP%
set GOARCH=%GOARCH_TMP%
set GOOS_TMP=
set GOARCH_TMP=
