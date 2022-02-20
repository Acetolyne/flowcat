#! /bin/bash

#Arch is passed in as linux-amd64, linux-i386, etc
arch=$1

cd ..
go install -v
go get gopkg.in/yaml.v2
go get github.com/Acetolyne/commentlex
#Test that building the main.go file matches the binary that is in bin, confirms the binary is the latest build
#//@todo push built binary only on PR (new workflow) to repo then remove diff checks below
#//@todo add more archs for binaries
case $arch in

  "linux-amd64")
    dirsep="/"
    env GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o bin/flowcat-linux-amd64/flowcat
    git branch
    git add bin/flowcat-linux-amd64/flowcat
    git commit -m "linux 386 auto build binary"

    #Everything up-to-date
    #diff flowcat bin/flowcat-$arch/flowcat
    #if [ `echo $?` -ne 0 ]; then echo "Binary file $1 is not the latest" && exit 1; fi
    ;;

  "linux-386")
    dirsep="/"
    env GOOS=linux GOARCH=386 GO111MODULE=auto go build -o bin/flowcat-linux-386/flowcat
    git status
    #diff flowcat bin/flowcat-$arch/flowcat
    #if [ `echo $?` -ne 0 ]; then echo "Binary file $1 is not the latest" && exit 1; fi
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
