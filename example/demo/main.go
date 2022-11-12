package main

import (
	"os"
	"strconv"

	muraylib "github.com/gabstv/microui-go-raylib"
	mu "github.com/gabstv/microui-go/microui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var flags uint32
	if ok, _ := strconv.ParseBool(os.Getenv("HIGHDPI")); ok {
		flags |= rl.FlagWindowHighdpi
	}
	if ok, _ := strconv.ParseBool(os.Getenv("RESIZE")); ok {
		flags |= rl.FlagWindowResizable
	}
	rl.SetConfigFlags(flags)
	rl.InitWindow(800, 450, "microui-go + raylib-go")
	rl.SetTargetFPS(60)

	ctx := mu.NewContext()
	muraylib.Setup(ctx)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		muraylib.UpdateInputs(ctx)
		ctx.Begin()
		mu.DrawDemoWindow(ctx)
		ctx.End()

		rl.ClearBackground(rl.Black)

		// rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
		ctx.Render()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
