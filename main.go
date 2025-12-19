package main

import (
	"game/source-code/global"
	"game/source-code/screens/battle"
	"game/source-code/screens/menu"
	"game/source-code/screens/world"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
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

	global.ThemesGUI = file.LoadText("data/gui/reusable-themes.xml")
	global.PopupDimGUI = file.LoadText("data/gui/reusable-popup-dim.xml")
	global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))
	global.ScreenMainMenu = screens.Add(menu.New(), true)
	global.ScreenWorld = screens.Add(world.New("data/worlds/test/map.tmx"), true)
	global.ScreenBattle = screens.Add(battle.New("data/battlegrounds/test/map.tmx"), true)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
			global.Project = tiled.NewProject(assets.LoadTiledProject("project.tiled-project"))
			screens.Reload()
		}
	}
}
