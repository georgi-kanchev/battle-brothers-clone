package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled/tilemap"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

var maps []*graphics.Sprite
var shapes []*geometry.Shape

func main() {
	var camera = graphics.NewCamera(1)
	reload()

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.DrawSprites(maps...)

		camera.DrawCircle(0, 0, 100, color.Red)

		for _, s := range shapes {
			var pts = s.CornerPoints()
			camera.DrawShapes(color.FadeOut(color.Red, 0.65), pts...)
			camera.DrawLinesPath(2, color.Red, pts...)
		}

		camera.DragAndZoom()
		if keyboard.IsKeyJustPressed(key.F5) {
			reload()
		}
	}
}

func reload() {
	assets.LoadTiledMap("maps/test/map.tmx")
	maps = tilemap.LayerSprites("maps/test/map.tmx", "Background", "")
	shapes = tilemap.LayerShapes("maps/test/map.tmx", "Zones", "")
}
