package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled/tilemap"
	"pure-game-kit/window"
)

func main() {
	var maps = reload()
	var camera = graphics.NewCamera(1)

	reload()

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.DrawSprites(maps...)

		camera.DragAndZoom()
		if keyboard.IsKeyPressedOnce(key.F5) {
			maps = reload()
		}
	}
}

func reload() []*graphics.Sprite {
	assets.LoadTiledMap("maps/test/map.tmx")
	var layer = tilemap.LayerSprites("maps/test/map", "Background", "")
	return layer
}
