//#include "raylib.h"
#include <stdio.h>
#include <stdlib.h>
#include <emscripten/emscripten.h>

const int screenWidth = 800;
const int screenHeight = 450;

//void UpdateDrawFrame(void);     // Update and Draw one frame
  
EMSCRIPTEN_KEEPALIVE
uint8_t* RenderPlay(int value1) {
    printf("digg = %d\n", value1);
    uint8_t* arr = malloc(sizeof(uint8_t)*5);
    arr[0] = 72;
    arr[1] = 101;
    arr[2] = 108;
    arr[3] = 108;
    arr[4] = 110;
  return arr;
}

int main(void)
{
    //RenderTexture2D game = LoadRenderTexture(screenWidth, screenHeight);
    //BeginDrawing();
    //    BeginTextureMode(game);
    //        ClearBackground(RAYWHITE);
    //        DrawText("Congrats! You created your first window!", 190, 200, 20, LIGHTGRAY);
    //    EndTextureMode();
    //EndDrawing();

    //Image frame = LoadImageFromTexture(game.texture);

    //for(int i=0;i<30;i++){

    //}
    return 0; 
}



//void UpdateDrawFrame(void){
//   BeginDrawing();
//
//   ClearBackground(RAYWHITE);
//
//   DrawText("Congrats! You created your first window!", 190, 200, 20, LIGHTGRAY);
//
//   EndDrawing();
//}
