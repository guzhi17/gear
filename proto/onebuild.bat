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

cd /d %~dp0

SET FILEPATH=%1

@rem check file exist
if not exist %FILEPATH% (
 echo not found: %FILEPATH%
 goto :eof
)

@rem get file name and path
for %%F in (%FILEPATH%) do (
echo %%~nxF %%~dpF
SET PROTOPATH=%%~dpF
SET FILE=%%~nxF
echo path: %PROTOPATH%
echo file: %FILE%
)


SET ROOTDIR=%cd%
set OUTDIR=%ROOTDIR%\\build

echo proto path: %PROTOPATH%
echo %ROOTDIR%


@rem PATH=%GOPATH%/bin;./protoc-3.6.1-win32/bin/;%PATH%
@rem PATH=%GOLIB%/bin;%GOLIB%/src/bin;%ROOTDIR%/protoc-3.14.0-win64/bin/;%PATH%
PATH=%ROOTDIR%/protoc-3.14.0-win64/bin/;%PATH%



echo outdir %OUTDIR%

protoc --version

echo build basic

call:BuildModuleS %FILE%

goto :eof

@rem /////////////////////////////////////////////////
:BuildAllProtoJs

set LAN=js

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%


protoc --%LAN%_out=import_style=commonjs,binary:%OUTDIR%/%LAN%/ %1


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
    call:BuildAllProto %%G %1
)
call:BuildJavaLiteProto %1
call:BuildAllProtoJs %1

cd /d %pwd%
echo enter %cd%

goto :eof
@rem /////////////////////////////////////////////////


:BuildAllProto

set LAN=%1

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%


protoc --%LAN%_out=%OUTDIR%/%LAN%/ %2


goto :eof

:BuildJavaLiteProto

set LAN=java

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%_lite


protoc --%LAN%_out=lite:%OUTDIR%/%LAN%_lite/ %1


goto :eof
@rem /////////////////////////////////////////////////
:BuildGoGrpc

set LAN=go

echo build %LAN%

echo mkdir: %OUTDIR%/%LAN%

mkdir %OUTDIR%\\%LAN%


protoc --%LAN%_out=grpc:%OUTDIR%/%LAN%/ %1


goto :eof

@rem /////////////////////////////////////////////////



@rem /////////////////////////////////////////////////

