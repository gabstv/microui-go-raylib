module github.com/gabstv/microui-go-raylib

go 1.19

require (
	github.com/gabstv/microui-go v0.0.0-20221112220741-854f4f28a28c
	github.com/gen2brain/raylib-go/raylib v0.0.0-00010101000000-000000000000
)

require golang.org/x/exp v0.0.0-20221111204811-129d8d6c17ab // indirect

replace github.com/gen2brain/raylib-go/raylib => ../raylib-go/raylib

replace github.com/gen2brain/raylib-go/raygui => ../raylib-go/raygui

replace github.com/gen2brain/raylib-go/easings => ../raylib-go/easings

replace github.com/gabstv/microui-go/microui => ../microui-go/microui
replace github.com/gabstv/microui-go => ../microui-go
