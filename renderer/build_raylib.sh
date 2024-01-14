export EMSDK_PATH="/home/magisterbrownie/Programs/emsdk"
export EMSCRIPTEN_PATH="$EMSDK_PATH/upstream/emscripten"
export CLANG_PATH="$EMSDK_PATH/upstream/bin"
export PYTHON_PATH="/usr/bin/python3"
export NODE_PATH="$EMSDK_PATH/node/12.9.1_64bit/bin"
#export PATH = $(shell printenv PATH):$(EMSDK_PATH):$(EMSCRIPTEN_PATH):$(CLANG_PATH):$(NODE_PATH):$(PYTHON_PATH)
export PATH=$(printenv PATH):$EMSDK_PATH:$EMSCRIPTEN_PATH:$CLANG_PATH:$NODE_PATH:$PYTHON_PATH
echo $PATH
cd raylib/src
make PLATFORM=PLATFORM_WEB -B

