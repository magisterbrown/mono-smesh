#emcc -o game.html main.c -Os -Wall ./raylib/src/libraylib.a -I. -I./raylib/src -L. -L./raylib/src -s USE_GLFW=3 --shell-file raylib/src/shell.html -DPLATFORM_WEB -sGL_ENABLE_GET_PROC_ADDRESS  -s EXTRA_EXPORTED_RUNTIME_METHODS='["ccall"]' -s EXPORTED_FUNCTIONS='["renderPlay"]'
emcc main.c -o game.js -s EXPORTED_RUNTIME_METHODS=['ccall']
