package world

import (
	"game/source-code/global"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	gfx "pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/time"
)

type World struct {
	worldPath, guiPath string
	tmap               *tiled.Map
	hud                *gui.GUI
	camera             *gfx.Camera

	parties []*Party
}

func New(worldPath, guiPath string) *World {
	var world = &World{worldPath: worldPath, guiPath: guiPath, camera: gfx.NewCamera(1)}
	world.parties = []*Party{NewParty(nil, 2250, 1530, true)}
	return world
}

//=================================================================

func (world *World) OnLoad() {
	world.tmap = tiled.NewMap(assets.LoadTiledMap(world.worldPath), global.Project)
	world.hud = gui.NewFromXML(file.LoadText(world.guiPath))
}
func (world *World) OnEnter() {}
func (world *World) OnUpdate() {
	world.camera.SetScreenAreaToWindow()
	world.tmap.Draw(world.camera)

	//=================================================================
	// parties
	for _, party := range world.parties {
		party.Update()
	}

	//=================================================================
	// gui
	if keyboard.IsKeyJustPressed(key.I) {
		var hidden = condition.If(world.hud.Field("inventory", field.Hidden) == "", "1", "")
		world.hud.SetField("popup-dim", field.Hidden, hidden)
		world.hud.SetField("inventory", field.Hidden, hidden)
		time.SetScale(condition.If(hidden == "1", float32(1), 0))
	} else if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
	}

	world.hud.UpdateAndDraw(world.camera)
}
func (world *World) OnExit() {}

//=================================================================
// private
