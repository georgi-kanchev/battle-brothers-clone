package battle

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
)

type Battle struct {
	path   string
	tmap   *tiled.Map
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	mapLayers []*tiled.Layer

	attackers, defenders []*unit.Unit
}

func New(mapPath string) *Battle {
	var battle = &Battle{path: mapPath, camera: graphics.NewCamera(1)}
	return battle
}

//=================================================================

func (battle *Battle) OnLoad() {
	battle.tmap = tiled.NewMap(assets.LoadTiledMap(battle.path), global.Project)
	battle.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	battle.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	battle.currentPopup = nil
	battle.mapLayers = battle.tmap.FindLayersBy(property.LayerClass, "BattleMap")
}
func (battle *Battle) OnEnter() {
	var cols = battle.tmap.Properties[property.MapColumns].(int)
	var rows = battle.tmap.Properties[property.MapRows].(int)
	var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
	var tileH = battle.tmap.Properties[property.MapTileHeight].(int)

	battle.camera.X, battle.camera.Y = float32(cols)/2*float32(tileW), float32(rows)/2*float32(tileH)
	battle.camera.Zoom = 0.8

	battle.attackers, battle.defenders = nil, nil
	for range 5 {
		battle.attackers = append(battle.attackers, unit.New(cols, rows, tileW, tileH))
		battle.defenders = append(battle.defenders, unit.New(cols, rows, tileW, tileH))
	}

	battle.spawnUnits(battle.attackers, false, "BattleSpawnsAttackers", tileW, tileH)
	battle.spawnUnits(battle.defenders, true, "BattleSpawnsDefenders", tileW, tileH)

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}

func (battle *Battle) spawnUnits(units []*unit.Unit, flip bool, layerClass string, tileW int, tileH int) {
	var spawns = battle.tmap.FindLayersBy(property.LayerClass, layerClass)[0].ExtractPoints()
	if len(units) > len(spawns) {
		return
	}
	for i, u := range units {
		var x = int(spawns[i][0] / float32(tileW))
		var y = int(spawns[i][1] / float32(tileH))
		u.Spawn(float32(x), float32(y), flip)
	}
}
func (battle *Battle) OnUpdate() {
	battle.camera.SetScreenAreaToWindow()

	if battle.currentPopup == nil {
		battle.camera.MouseDragAndZoomSmooth()
		battle.camera.Zoom = number.Limit(battle.camera.Zoom, 0.5, 10)
	}
	for _, l := range battle.mapLayers {
		l.Draw(battle.camera)
	}

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
	var allUnits = collection.Join(battle.attackers, battle.defenders)

	for _, unit := range allUnits {
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
