package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled/tilemap"
	"pure-game-kit/utility/collection"
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
	assets.LoadTextures("maps/ground.png")
	assets.LoadTiledTileset("maps/ground.tsx")
	assets.LoadTiledWorld("maps/maps.world")

	var map1 = tilemap.LayerSprites("maps/world1", "Tile Layer 1", "")
	var map2 = tilemap.LayerSprites("maps/world2", "Tile Layer 1", "")
	return collection.Join(map1, map2)
}
