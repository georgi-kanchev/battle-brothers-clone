package world

import (
	"game/source-code/global"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
)

type World struct {
	path   string
	tmap   *tiled.Map
	camera *graphics.Camera

	hud, inventory, settlement, currentPopup *gui.GUI

	parties []*Party
}

func New(path string) *World {
	var world = &World{path: path, camera: graphics.NewCamera(1)}
	world.parties = []*Party{NewParty(nil, 2250, 1530, true)}
	return world
}

//=================================================================

func (world *World) OnLoad() {
	world.tmap = tiled.NewMap(assets.LoadTiledMap(world.path), global.Project)
	world.hud = gui.NewFromXMLs(file.LoadText("data/gui/world-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	world.inventory = gui.NewFromXMLs(file.LoadText("data/gui/world-inventory.xml"), global.ThemesGUI)
	world.settlement = gui.NewFromXMLs(file.LoadText("data/gui/world-settlement.xml"), global.ThemesGUI)
	world.currentPopup = nil
}
func (world *World) OnEnter() {
}
func (world *World) OnUpdate() {
	world.camera.SetScreenAreaToWindow()
	world.tmap.Draw(world.camera)

	for _, party := range world.parties {
		party.Update()
	}

	world.hud.UpdateAndDraw(world.camera)
	if world.currentPopup != nil {
		world.currentPopup.UpdateAndDraw(world.camera)
	}

	world.handleInput()
}

func (world *World) OnExit() {
}

//=================================================================
// private

func (world *World) handleInput() {
	if keyboard.IsKeyJustPressed(key.I) {
		world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.inventory)
	} else if keyboard.IsKeyJustPressed(key.S) {
		world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.settlement)
	} else if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
	} else if keyboard.IsKeyJustPressed(key.Escape) {
		if world.currentPopup == nil {
			screens.Enter(global.ScreenMainMenu, false)
		} else {
			world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.currentPopup)
		}
	} else if world.settlement.IsButtonJustClicked("settlement-exit-btn", world.camera) {
		world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.settlement)
	}
}
