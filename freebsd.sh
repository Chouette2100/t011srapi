#/bin/sh
if [ $# -ne 1 ]; then
  echo "指定された引数は$#個です。" 1>&2
  echo "実行するには1個の引数が必要です。" 1>&2
  echo "usage: $0 source_filename[.go]"
  exit 1
fi
filename=$1
filename=${filename%.*}
export GOOS_TMP=$GOOS
export GOARCH_TMP=$GOARCH
export GOOS=freebsd
export GOARCH=amd64
go build -o $filename.freebsd $filename.go
export GOOS=$GOOS_TMP
export GOARCH=$GOARCH_TMP
unset GOOS_TMP
unset GOARCH_TMP
