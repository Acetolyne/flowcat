#! /bin/bash

#Arch is passed in as linux-amd64, linux-i386, etc
arch=$1
#Test that building the main.go file matches the binary that is in bin, confirms the binary is the latest build
case $arch in

  "linux-amd64")
    dirsep="/"
    cd ..
    which go
    go mod init main
    go get gopkg.in/yaml.v2
    go install -v
    env GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o flowcat
    #diff flowcat bin/flowcat-$arch/flowcat
    #if [ `echo $?` -ne 0 ]; then echo "Binary file $1 is not the latest" && exit 1; fi
    ;;

  "linux-386")
    dirsep="/"
    cd ..
    go get gopkg.in/yaml.v2
    env GOOS=linux GOARCH=386 GO111MODULE=auto go build -o flowcat
    diff flowcat bin/flowcat-$arch/flowcat
    if [ `echo $?` -ne 0 ]; then echo "Binary file $1 is not the latest" && exit 1; fi
    ;;

  *)
    echo "Unable to build for $arch" && exit 1
    ;;
esac

#Start functionality tests for all binaries
####TEMP STUFF
# cd ..
# ls
# echo $GITHUB_PATH
# mkdir tests/bin
# cp bin/flowcat-linux-amd64/flowcat tests/bin/flowcat

#Setup directory structure for tests

sudo mkdir -p tests/tmp
sudo mkdir -p tests/tmp2
sudo mkdir -p tests/bin
sudo mv flowcat /bin/flowcat
#echo "$GITHUB_WORKSPACE/tests/bin" >> $GITHUB_PATH

cd tests/tmp
sudo cp ../assets/__testfile__ .
cd ../tmp2
sudo cp ../assets/__testfile__ .
sudo cp ../assets/.test .
sudo cp ../assets/regular .
cd ../tmp
pwd
#CanRun Test
res=$(flowcat)
echo $res | grep -q '__testfile__ test file after'
if [ $(echo $?) -ne 0 ]; then
  echo "CanRun Failed"
  exit 1
fi
#CanOutputLinenums
res=$(flowcat -l)
echo $res | grep -q '__testfile__ 1) test file 3) after'
if [ $(echo $?) -ne 0 ]; then
  echo "CanOutputLinenums Failed"
  exit 1
fi
#CanSpecifyMatch
res=$(flowcat -m "#@todo")
echo $res | grep -q '__testfile__ test 2'
if [ $(echo $?) -ne 0 ]; then
  echo "CanSpecifyMatch Failed"
  exit 1
fi
#CanCreateOutputFile
flowcat -o todo
res=$(cat todo)
echo $res | grep -q '__testfile__ test file after'
if [ $(echo $?) -ne 0 ]; then
  echo "CanCreateOutputFile Failed"
  exit 1
fi
#CanDisplayHelp
res=$(flowcat -h)
echo $res | grep -q 'Options for Flowcat'
if [ $(echo $?) -ne 0 ]; then
  echo "CanDisplayHelp Failed"
  exit 1
fi
#CanUseSettingsFile
echo "stat1"
ls -la ../assets/
sudo cp ../assets/.flowcat .
res=$(flowcat)
echo $res | grep -q '__testfile__ 1) test file 3) after'
if [ $(echo $?) -ne 0 ]; then
  echo "CanUseSettingsFile Failed"
  exit 1
fi
sudo cp ../assets/.flowcat1 ./.flowcat
res=$(flowcat)
echo $res | grep -q '__testfile__ test 2'
if [ $(echo $?) -ne 0 ]; then
  echo "CanUseSettingsFile Failed"
  exit 1
fi
#CanSpecifyPath
res=$(flowcat -f ../tmp2/)
echo $res | grep -q '../tmp2/.test test file after ../tmp2/__testfile__ test file after ../tmp2/regular regular test with exclude'
if [ $(echo $?) -ne 0 ]; then
  echo "CanSpecifyPath Failed"
  exit 1
fi
#CanUsePathSettings
sudo cp ../assets/.flowcat ../tmp2/.flowcat
cd ..
res=$(flowcat -f tmp2/)
echo $res | grep -q 'tmp2/.test 1) test file 3) after tmp2/__testfile__ 1) test file 3) after tmp2/regular 1) regular test 3) with exclude'
if [ $(echo $?) -ne 0 ]; then
  echo "CanUsePathSettings Failed"
  exit 1
fi
