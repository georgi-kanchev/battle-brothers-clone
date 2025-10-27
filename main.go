package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled/tilemap"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func main() {
	var maps = reload()
	var camera = graphics.NewCamera(1)

	reload()

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.DrawSprites(maps...)

		camera.DrawGrid(2, 128, 128, color.Darken(color.Gray, 0.5))
		camera.DragAndZoom()

		if keyboard.IsKeyPressedOnce(key.F5) {
			maps = reload()
		}
	}
}

func reload() []*graphics.Sprite {
	assets.LoadTiledMap("maps/world.tmx")
	return tilemap.LayerSprites("maps/world", "Tile Layer 1", "")
}
