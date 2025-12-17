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
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
)

type Battle struct {
	path   string
	tmap   *tiled.Map
	camera *gfx.Camera

	hud, currentPopup, loot *gui.GUI

	attackers, defenders []*unit.Unit
}

func New(mapPath string) *Battle {
	var battle = &Battle{path: mapPath, camera: gfx.NewCamera(1)}
	return battle
}

//=================================================================

func (battle *Battle) OnLoad() {
	battle.tmap = tiled.NewMap(assets.LoadTiledMap(battle.path), global.Project)
	battle.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	battle.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	battle.currentPopup = nil

	var cols = battle.tmap.Properties[property.MapColumns].(int)
	var rows = battle.tmap.Properties[property.MapRows].(int)
	var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
	var tileH = battle.tmap.Properties[property.MapTileHeight].(int)
	battle.camera.X, battle.camera.Y = float32(cols)/2*float32(tileW), float32(rows)/2*float32(tileH)
}
func (battle *Battle) OnEnter() {
}
func (battle *Battle) OnUpdate() {
	battle.camera.SetScreenAreaToWindow()
	condition.CallIf(battle.currentPopup == nil, battle.camera.MouseDragAndZoomSmooth)
	battle.tmap.Draw(battle.camera)

	//=================================================================
	// units
	for _, unit := range battle.ySortUnits() {
		unit.Draw(battle.camera)
	}

	//=================================================================
	// gui

	if keyboard.IsKeyJustPressed(key.W) {
		screens.Enter(global.ScreenWorld, false)
	}

	if keyboard.IsKeyJustPressed(key.L) {
		battle.currentPopup = global.TogglePopup(battle.hud, battle.currentPopup, battle.loot)
	}

	battle.hud.UpdateAndDraw(battle.camera)

	if battle.currentPopup != nil {
		battle.currentPopup.UpdateAndDraw(battle.camera)
	}
}
func (battle *Battle) OnExit() {
}

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
