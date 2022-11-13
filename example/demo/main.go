package main

import (
	"image/color"
	"os"
	"strconv"

	muraylib "github.com/gabstv/microui-go-raylib"
	"github.com/gabstv/microui-go/demo"
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
	rl.InitWindow(800, 500, "microui-go + raylib-go")
	rl.SetTargetFPS(60)

	ctx := mu.NewContext()
	muraylib.Setup(ctx)

	for !rl.WindowShouldClose() {
		muraylib.UpdateInputs(ctx)
		ctx.Begin()
		demo.DemoWindow(ctx)
		demo.LogWindow(ctx)
		demo.StyleWindow(ctx)
		ctx.End()

		rl.BeginDrawing()

		// get the color from the demo window (because of the sliders)
		bgc := demo.BackgroundColor()
		cc := color.RGBA{
			R: bgc.R,
			G: bgc.G,
			B: bgc.B,
			A: bgc.A,
		}

		rl.ClearBackground(cc)

		ctx.Render()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
