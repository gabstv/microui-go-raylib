package muraylib

import (
	_ "embed"
	"sync"

	"github.com/gabstv/microui-go/microui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var initOnce sync.Once

//go:embed default_atlas.png
var defaultAtlasPNG []byte

var defaultAtlas *rl.Image
var DefaultAtlasTexture rl.Texture2D

var (
	iconClose     = microui.NewRect(0, 0, 16, 16)
	iconCheck     = microui.NewRect(16, 0, 16, 16)
	iconCollapsed = microui.NewRect(32, 0, 16, 16)
	iconExpanded  = microui.NewRect(48, 0, 16, 16)
	atlasWhite    = microui.NewRect(2, 18, 3, 3)

	DefaultAtlasRects = []microui.Rect{
		{},
		iconClose,
		iconCheck,
		iconCollapsed,
		iconExpanded,
		atlasWhite,
	}
)

// raylib default font size
const defaultFontSize = 10

func Setup(ctx *microui.Context) {
	initOnce.Do(atlasSetup)
	ctx.SetRenderCommand(RenderCommand)
}

func RenderCommand(cmd *microui.Command) {
	switch cmd.Type() {
	case microui.CommandText:
		renderText(cmd.Text())
	case microui.CommandRect:
		renderRect(cmd.Rect())
	case microui.CommandIcon:
		renderIcon(cmd.Icon())
	case microui.CommandClip:
		renderClip(cmd.Clip())
	}
}

func renderText(cmd microui.TextCommand) {
	pos := cmd.Pos()
	fnt := cmd.Font()
	c := cmd.Color()
	txt := cmd.Text()
	if uintptr(fnt) == 0 {
		rl.DrawText(txt, int32(pos.X), int32(pos.Y), defaultFontSize, rl.NewColor(c.R, c.G, c.B, c.A))
		return
	}
	//TODO: this
}

func renderRect(cmd microui.RectCommand) {
	rect := cmd.Rect()
	c := cmd.Color()
	rl.DrawRectangle(rect.X, rect.Y, rect.W, rect.H, rl.NewColor(c.R, c.G, c.B, c.A))
}

func renderIcon(cmd microui.IconCommand) {
	rect := DefaultAtlasRects[int(cmd.ID())]
	x := cmd.Rect().X + (cmd.Rect().W-rect.W)/2
	y := cmd.Rect().Y + (cmd.Rect().H-rect.H)/2
	renderAtlasTexture(rect, microui.Vec2{X: x, Y: y}, cmd.Color())
}

func renderClip(cmd microui.ClipCommand) {
	rl.EndScissorMode()
	// AMENDED TO APPLY SCISSOR MODE CORRECTLY
	rect := cmd.Rect()
	rl.BeginScissorMode(rect.X, int32(rl.GetScreenHeight())-(rect.Y+rect.H), rect.W, rect.H)
}

func renderAtlasTexture(rect microui.Rect, pos microui.Vec2, color microui.Color) {
	source := rl.NewRectangle(float32(rect.X), float32(rect.Y), float32(rect.W), float32(rect.H))
	position := rl.NewVector2(float32(pos.X), float32(pos.Y))
	rl.DrawTextureRec(DefaultAtlasTexture, source, position, rl.NewColor(color.R, color.G, color.B, color.A))
}

func atlasSetup() {
	defaultAtlas = rl.LoadImageFromMemory(".png", defaultAtlasPNG, int32(len(defaultAtlasPNG)))
	if defaultAtlas == nil {
		panic("failed to load default atlas")
	}
	println("default atlas loaded")
	DefaultAtlasTexture = rl.LoadTextureFromImage(defaultAtlas)
	// DefaultAtlasTexture = rl.LoadTexture("default_atlas.png")
}

func UpdateInputs(ctx *microui.Context) {
	// while (SDL_PollEvent(&e)) {
	// 	switch (e.type) {
	// 	  case SDL_QUIT: exit(EXIT_SUCCESS); break;
	// 	  case SDL_MOUSEMOTION: mu_input_mousemove(ctx, e.motion.x, e.motion.y); break;
	// 	  case SDL_MOUSEWHEEL: mu_input_scroll(ctx, 0, e.wheel.y * -30); break;
	// 	  case SDL_TEXTINPUT: mu_input_text(ctx, e.text.text); break;

	// 	  case SDL_MOUSEBUTTONDOWN:
	// 	  case SDL_MOUSEBUTTONUP: {
	// 		int b = button_map[e.button.button & 0xff];
	// 		if (b && e.type == SDL_MOUSEBUTTONDOWN) { mu_input_mousedown(ctx, e.button.x, e.button.y, b); }
	// 		if (b && e.type ==   SDL_MOUSEBUTTONUP) { mu_input_mouseup(ctx, e.button.x, e.button.y, b);   }
	// 		break;
	// 	  }

	// 	  case SDL_KEYDOWN:
	// 	  case SDL_KEYUP: {
	// 		int c = key_map[e.key.keysym.sym & 0xff];
	// 		if (c && e.type == SDL_KEYDOWN) { mu_input_keydown(ctx, c); }
	// 		if (c && e.type ==   SDL_KEYUP) { mu_input_keyup(ctx, c);   }
	// 		break;
	// 	  }
	// 	}
	//   }
	md := rl.GetMousePosition()
	ctx.InputMouseMove(int32(md.X), int32(md.Y))

	mw := rl.GetMouseWheelMoveV()
	if mw.X != 0 || mw.Y != 0 {
		ctx.InputScroll(int32(mw.X), int32(mw.Y))
	}
	var mbtns microui.MouseButton
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mbtns |= microui.MouseLeft
	}
	if rl.IsMouseButtonPressed(rl.MouseMiddleButton) {
		mbtns |= microui.MouseMiddle
	}
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		mbtns |= microui.MouseRight
	}
	if mbtns != 0 {
		mp := rl.GetMousePosition()
		ctx.InputMouseDown(int32(mp.X), int32(mp.Y), mbtns)
	}
	var mbtnsUp microui.MouseButton
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mbtnsUp |= microui.MouseLeft
	}
	if rl.IsMouseButtonReleased(rl.MouseMiddleButton) {
		mbtnsUp |= microui.MouseMiddle
	}
	if rl.IsMouseButtonReleased(rl.MouseRightButton) {
		mbtnsUp |= microui.MouseRight
	}
	if mbtnsUp != 0 {
		mp := rl.GetMousePosition()
		ctx.InputMouseUp(int32(mp.X), int32(mp.Y), mbtnsUp)
	}
}

func init() {
	// runtime.LockOSThread()
	microui.DefaultGetTextWidth = func(font microui.Font, text string) int32 {
		if uintptr(font) == 0 {
			return rl.MeasureText(text, defaultFontSize)
		}
		// TODO: custom fonts
		// rl.Font
		// rl.GetFontDefault()
		// TODO
		return 1
	}
	microui.DefaultGetTextHeight = func(font microui.Font) int32 {
		if uintptr(font) == 0 {
			return rl.GetFontDefault().BaseSize
		}
		//TODO: custom fonts
		return 1
	}
}
