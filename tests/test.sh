#! /bin/bash

#Arch is passed in as linux-amd64, linux-i386, etc
arch=$1
#Test that building the main.go file matches the binary that is in bin, confirms the binary is the latest build
# case $arch in

#   "linux-amd64")
#     dirsep="/"
#     cd ..
#     env GOOS=linux GOARCH=amd64 go build -o flowcat
#     ls -la
#     diff flowcat bin/flowcat-$arch/flowcat
#     if [[ `echo $?` -ne 0 ]]; then echo "Binary file $1 is not the latest" && exit 1; fi
#     ;;

#   "linux-386")
#     dirsep="/"
#     cd ..
#     env GOOS=linux GOARCH=386 go build -o flowcat
#     diff flowcat bin/flowcat-$arch/flowcat
#     if [[ `echo $?` -ne 0 ]]; then echo "Binary file $1 is not the latest" && exit 1; fi
#     ;;

#   *)
#     echo "Unable to build for $arch" && exit 1
#     ;;
# esac

#Start functionality tests for all binaries

#Setup directory structure for tests
mkdir tests/tmp
mkdir tests/tmp2
mkdir tests/bin
mv flowcat tests/bin/flowcat
echo "::add-path::$GITHUB_WORKSPACE/tests/bin"
cd tests/tmp
cp ../assets/__testfile__ .
cd ../tmp2
cp ../assets/__testfile__ .
cp ../assets/.test .
cp ../assets/regular .
cd ../tmp
#CanRun Test
res=$(flowcat)
echo $res | grep -q '__testfile__ test file after'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanRun Failed"
  exit 1
fi
#CanOutputLinenums
res=$(flowcat -l)
echo $res | grep -q '__testfile__ 1) test file 3) after'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanOutputLinenums Failed"
  exit 1
fi
#CanSpecifyMatch
res=$(flowcat -m "#@todo")
echo $res | grep -q '__testfile__ test 2'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanSpecifyMatch Failed"
  exit 1
fi
#CanCreateOutputFile
flowcat -o todo > /dev/null
res=$(cat todo)
echo $res | grep -q '__testfile__ test file after'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanCreateOutputFile Failed"
  exit 1
fi
#CanDisplayHelp
res=$(flowcat -h)
echo $res | grep -q 'Options for Flowcat'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanDisplayHelp Failed"
  exit 1
fi
#CanUseSettingsFile
cp ../assets/.flowcat .
res=$(flowcat)
echo $res | grep -q '__testfile__ 1) test file 3) after'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanUseSettingsFile Failed"
  exit 1
fi
cp ../assets/.flowcat1 ./.flowcat
res=$(flowcat)
echo $res | grep -q '__testfile__ test 2'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanUseSettingsFile Failed"
  exit 1
fi
#CanSpecifyPath
res=$(flowcat -f ../tmp2/)
echo $res | grep -q '../tmp2/.test test file after ../tmp2/__testfile__ test file after ../tmp2/regular regular test with exclude'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanSpecifyPath Failed"
  exit 1
fi
#CanUsePathSettings
cp ../assets/.flowcat ../tmp2/.flowcat
cd ..
res=$(flowcat -f tmp2/)
echo $res | grep -q 'tmp2/.test 1) test file 3) after tmp2/__testfile__ 1) test file 3) after tmp2/regular 1) regular test 3) with exclude'
if [[ $(echo $?) -ne 0 ]]; then
  echo "CanUsePathSettings Failed"
  exit 1
fi
