package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/window"
)

func main() {
	var camera = graphics.NewCamera(1)
	var assetId = assets.LoadTiledMap("maps/test/map.tmx")
	var projectId = assets.LoadTiledProject("maps/map.tiled-project")
	var project = tiled.NewProject(projectId)
	var mapp = tiled.NewMap(assetId, project)

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.MouseDragAndZoomSmooth()

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
		}

		mapp.Draw(camera)
	}
}
