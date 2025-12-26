package battle

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
)

type Battle struct {
	path   string
	Tmap   *tiled.Map
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	attackers, defenders []*unit.Unit
}

func New(mapPath string) *Battle {
	var battle = &Battle{path: mapPath, camera: graphics.NewCamera(1)}
	return battle
}

//=================================================================

func (battle *Battle) OnLoad() {
	battle.Tmap = tiled.NewMap(assets.LoadTiledMap(battle.path), global.Project)
	battle.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	battle.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	battle.currentPopup = nil

	var cols = battle.Tmap.Properties[property.MapColumns].(int)
	var rows = battle.Tmap.Properties[property.MapRows].(int)
	var tileW = battle.Tmap.Properties[property.MapTileWidth].(int)
	var tileH = battle.Tmap.Properties[property.MapTileHeight].(int)
	battle.camera.X, battle.camera.Y = float32(cols)/2*float32(tileW), float32(rows)/2*float32(tileH)

	battle.attackers = append(battle.attackers, unit.New(cols, rows, tileW, tileH))

	var all = collection.Join(battle.attackers, battle.defenders)
	for _, v := range all {
		v.Load()
	}
}
func (battle *Battle) OnEnter() {
}
func (battle *Battle) OnUpdate() {
	battle.camera.SetScreenAreaToWindow()
	condition.CallIf(battle.currentPopup == nil, battle.camera.MouseDragAndZoomSmooth)
	battle.Tmap.Draw(battle.camera)

	for _, unit := range battle.ySortUnits() {
		unit.Draw(battle.camera)
	}

	battle.hud.UpdateAndDraw(battle.camera)
	if battle.currentPopup != nil {
		battle.currentPopup.UpdateAndDraw(battle.camera)
	}

	battle.handleInput()
}

func (battle *Battle) OnExit() {
}

//=================================================================
// private

func (battle *Battle) handleInput() {
	if keyboard.IsKeyJustPressed(key.Escape) {
		screens.Enter(global.ScreenWorld, false)
	} else if keyboard.IsKeyJustPressed(key.L) {
		battle.currentPopup = global.TogglePopup(battle.hud, battle.currentPopup, battle.loot)
	}
}
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
