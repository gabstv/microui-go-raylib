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

var DefaultScrollMultiplier = rl.NewVector2(1, -30)

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
	ctx.SetBeginRender(func() {
		rl.BeginScissorMode(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()))
	})
	ctx.SetEndRender(func() {
		rl.EndScissorMode()
	})
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
	rect := cmd.Rect()
	rl.BeginScissorMode(rect.X, rect.Y, rect.W, rect.H)
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

var charbuf = make([]rune, 1024)

func UpdateInputs(ctx *microui.Context) {
	md := rl.GetMousePosition()
	ctx.InputMouseMove(int32(md.X), int32(md.Y))

	mw := rl.GetMouseWheelMoveV()
	if mw.X != 0 || mw.Y != 0 {
		ctx.InputScroll(int32(mw.X*DefaultScrollMultiplier.X), int32(mw.Y*DefaultScrollMultiplier.Y))
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
	nchars := 0
	lastc := rl.GetCharPressed()
	for lastc != 0 {
		charbuf[nchars] = lastc
		nchars++
		lastc = rl.GetCharPressed()
	}
	if nchars > 0 {
		ctx.InputText(string(charbuf[:nchars]))
	}

	if rl.IsKeyPressed(rl.KeyLeftControl) || rl.IsKeyPressed(rl.KeyRightControl) {
		ctx.InputKeyDown(microui.KeyCtrl)
	}
	if rl.IsKeyPressed(rl.KeyLeftShift) || rl.IsKeyPressed(rl.KeyRightShift) {
		ctx.InputKeyDown(microui.KeyShift)
	}
	if rl.IsKeyPressed(rl.KeyLeftAlt) || rl.IsKeyPressed(rl.KeyRightAlt) {
		ctx.InputKeyDown(microui.KeyAlt)
	}
	if rl.IsKeyPressed(rl.KeyBackspace) {
		ctx.InputKeyDown(microui.KeyBackspace)
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		ctx.InputKeyDown(microui.KeyReturn)
	}

	if rl.IsKeyReleased(rl.KeyLeftControl) || rl.IsKeyReleased(rl.KeyRightControl) {
		ctx.InputKeyUp(microui.KeyCtrl)
	}
	if rl.IsKeyReleased(rl.KeyLeftShift) || rl.IsKeyReleased(rl.KeyRightShift) {
		ctx.InputKeyUp(microui.KeyShift)
	}
	if rl.IsKeyReleased(rl.KeyLeftAlt) || rl.IsKeyReleased(rl.KeyRightAlt) {
		ctx.InputKeyUp(microui.KeyAlt)
	}
	if rl.IsKeyReleased(rl.KeyBackspace) {
		ctx.InputKeyUp(microui.KeyBackspace)
	}
	if rl.IsKeyReleased(rl.KeyEnter) {
		ctx.InputKeyUp(microui.KeyReturn)
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
