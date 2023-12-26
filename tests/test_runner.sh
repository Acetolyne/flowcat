#! /bin/bash

#Arch is passed in as linux-amd64, linux-i386, etc
arch=$1

cd ..
go install -v
go get gopkg.in/yaml.v2
#//@todo add more archs for binaries
case $arch in

  "linux-amd64")
    dirsep="/"
    env GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o bin/flowcat-linux-amd64/flowcat
    ;;

  "linux-386")
    dirsep="/"
    env GOOS=linux GOARCH=386 GO111MODULE=auto go build -o bin/flowcat-linux-386/flowcat
    ;;

  "darwin-arm64")
    dirsep="/"
    env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build -o bin/flowcat-darwin-arm64/flowcat
    ;;

  *)
    echo "Unable to build for $arch" && exit 1
    ;;
esac

#Start functionality tests for all binaries

res=$(go test -v)
if [ $(echo $?) -ne 0 ]; then
  echo "Tests Failed"
  exit 1
fi
echo "PASSED"
