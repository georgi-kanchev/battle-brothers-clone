package main

import (
	"game/code/global"
	"game/code/screens/battle"
	"game/code/screens/loading"
	"game/code/screens/menu"
	"game/code/screens/world"
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
	assets.LoadDefaultFont()
	global.LoadOptions()
	global.ApplyOptions()

	window.KeepOpen()
	global.ScreenLoading = screens.Add(loading.New(), true)

	loading.Show("Loading:\nReusable GUI...")
	global.ThemesGUI = file.LoadText("data/gui/reusable-themes.xml")
	global.PopupNarrowGUI = file.LoadText("data/gui/reusable-popup-narrow.xml")
	global.PopupWideGUI = file.LoadText("data/gui/reusable-popup-wide.xml")
	global.DimGUI = file.LoadText("data/gui/reusable-popup-dim.xml")
	global.XBtnGUI = file.LoadText("data/gui/reusable-popup-x-button.xml")
	global.TitleGUI = file.LoadText("data/gui/reusable-popup-title.xml")

	loading.Show("Loading:\nTiled project...")
	global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))

	loading.Show("Loading:\nScreens...")
	global.ScreenMainMenu = screens.Add(menu.New(), true)
	global.ScreenWorld = screens.Add(world.New("data/worlds/test/map.tmx"), true)
	global.ScreenBattle = screens.Add(battle.New("data/battlegrounds/test/map.tmx"), true)

	screens.Enter(global.ScreenMainMenu, false)
	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			var prevScreen = screens.CurrentId()
			global.LoadOptions()
			global.ApplyOptions()
			assets.ReloadAll()
			global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))
			screens.Reload()
			screens.Enter(prevScreen, false)
		}
	}
}
