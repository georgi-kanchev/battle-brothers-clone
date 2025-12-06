package battle

import (
	"game/source-code/global"
	"game/source-code/unit"
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
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
)

type Battle struct {
	mapPath, guiPath string
	tmap             *tiled.Map
	hud              *gui.GUI
	camera           *gfx.Camera

	attackers, defenders []*unit.Unit
}

func New(mapPath, guiPath string) *Battle {
	var battle = &Battle{mapPath: mapPath, guiPath: guiPath, camera: gfx.NewCamera(1)}
	return battle
}

//=================================================================

func (battle *Battle) OnLoad() {
	battle.tmap = tiled.NewMap(assets.LoadTiledMap(battle.mapPath), global.Project)
	battle.hud = gui.NewFromXML(file.LoadText(battle.guiPath))

	var cols = battle.tmap.Properties[property.MapColumns].(int)
	var rows = battle.tmap.Properties[property.MapRows].(int)
	var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
	var tileH = battle.tmap.Properties[property.MapTileHeight].(int)
	battle.camera.X, battle.camera.Y = float32(cols)/2*float32(tileW), float32(rows)/2*float32(tileH)
}
func (battle *Battle) OnEnter() {}
func (battle *Battle) OnUpdate() {
	battle.camera.SetScreenAreaToWindow()
	battle.camera.MouseDragAndZoomSmooth()
	battle.tmap.Draw(battle.camera)

	//=================================================================
	// units
	for _, unit := range battle.ySortUnits() {
		unit.Draw(battle.camera)
	}

	//=================================================================
	// gui
	battle.hud.UpdateAndDraw(battle.camera)

	if keyboard.IsKeyJustPressed(key.W) {
		screens.Enter(global.ScreenWorld, false)
	}

	if keyboard.IsKeyJustPressed(key.L) {
		var hidden = condition.If(battle.hud.Field("loot", field.Hidden) == "", "1", "")
		battle.hud.SetField("popup-dim", field.Hidden, hidden)
		battle.hud.SetField("loot", field.Hidden, hidden)
	}
}
func (battle *Battle) OnExit() {}

//=================================================================
// private

func (battle *Battle) ySortUnits() []*unit.Unit {
	var ySorted = make(map[float32][]*unit.Unit, len(battle.attackers)+len(battle.defenders))

	for _, unit := range battle.attackers {
		var _, y = unit.Position()
		ySorted[y] = append(ySorted[y], unit)
	}

	var keys = collection.MapKeys(ySorted)
	var result = make([]*unit.Unit, 0, len(ySorted))

	collection.SortNumbers(keys...)
	for _, key := range keys {
		result = append(result, ySorted[key]...)
	}

	return result
}
