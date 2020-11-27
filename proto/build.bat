@echo off
SET GOPATHCopy=%GOPATH%
SET GOLIB = %GOPATH%
:WHILE
IF "%GOPATHCopy%" == "" GOTO WEND
FOR /F "delims=;" %%I IN ("%GOPATHCopy%") DO (
SET GOLIB=%%~I
GOTO :WEND
)
:WEND

echo gopath: %GOPATHCopy%
echo golib: %GOLIB%


@rem SET EXCLUDES=(pb_moment.proto)

cd /d %~dp0


SET ROOTDIR=%cd%
echo %ROOTDIR%

@rem PATH=%GOPATH%/bin;./protoc-3.6.1-win32/bin/;%PATH%
PATH=%ROOTDIR%/protoc-3.14.0-win64/bin/;%PATH%

set PROTOPATH=./proto
set OUTDIR=%ROOTDIR%\\build

echo outdir %OUTDIR%

protoc --version

echo build basic

call:BuildModuleS

goto :eof

@rem /////////////////////////////////////////////////
:BuildAllProtoJs

set LAN=js

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%

for /f %%G in ('dir /b *.proto') do (
    protoc --%LAN%_out=import_style=commonjs,binary:%OUTDIR%/%LAN%/ %%G
)

echo "go js browserify"

cd %OUTDIR%\\%LAN%
echo %cd%

call browserify exports.js -o pbbat.js


goto :eof


@rem /////////////////////////////////////////////////
:BuildModuleS

@rem protoc --python_out=./python/ ./%PROTOPATH%/*.proto
@rem protoc --java_out=./java/ ./%PROTOPATH%/*.proto
@rem protoc --objc_out=./oc/ ./%PROTOPATH%/*.proto

set pwd=%cd%
cd %PROTOPATH%

echo %cd%

@rem go get github.com/golang/protobuf/protoc-gen-go
@rem,dart
for %%G in (python,java,objc,cpp,go) do (
    call:BuildAllProto %%G
)
call:BuildJavaLiteProto
call:BuildAllProtoJs

cd /d %pwd%
echo enter %cd%

goto :eof
@rem /////////////////////////////////////////////////


:BuildAllProto

set LAN=%1

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%

for /f %%G in ('dir /b *.proto') do (
    protoc --%LAN%_out=%OUTDIR%/%LAN%/ %%G
)

goto :eof

:BuildJavaLiteProto

set LAN=java

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%_lite

for /f %%G in ('dir /b *.proto') do (
    protoc --%LAN%_out=lite:%OUTDIR%/%LAN%_lite/ %%G
)

goto :eof
@rem /////////////////////////////////////////////////


@rem /////////////////////////////////////////////////



@rem /////////////////////////////////////////////////