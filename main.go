package main

import (
	"game/source-code/global"
	"game/source-code/screens/battle"
	"game/source-code/screens/world"
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/window"
)

func main() {
	window.Title = "Battle Brothers Clone"

	assets.LoadDefaultAtlasIcons()
	assets.LoadDefaultFont()
	assets.LoadDefaultTexture()

	global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))
	// global.ScreenMenu = screens.Add(nil, false)
	global.ScreenWorld = screens.Add(world.New("data/worlds/test/map.tmx", "data/gui/world.xml"), true)
	global.ScreenBattle = screens.Add(battle.New("data/battlegrounds/test/map.tmx", "data/gui/battle.xml"), true)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
			global.Project = tiled.NewProject(assets.LoadTiledProject("project.tiled-project"))
			screens.Reload()
		}
	}
}
