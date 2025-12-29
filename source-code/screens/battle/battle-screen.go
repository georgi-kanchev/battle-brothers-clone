package battle

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
)

type BattleScreen struct {
	path   string
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	tmap  *tiled.Map
	tiles []*graphics.Sprite

	units []*unit.Unit

	turnManager *turnManager
}

func New(mapPath string) *BattleScreen {
	var battle = &BattleScreen{path: mapPath, camera: graphics.NewCamera(1), turnManager: newTurnManager()}
	return battle
}

//=================================================================

func (b *BattleScreen) OnLoad() {
	b.tmap = tiled.NewMap(assets.LoadTiledMap(b.path), global.Project)
	b.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	b.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	b.currentPopup = nil

	var layers = b.tmap.FindLayersBy(property.LayerClass, "BattleMap")
	for _, l := range layers {
		b.tiles = append(b.tiles, l.ExtractSprites()...)
	}

	// assets.SetTextureSmoothness("art/Battlegrounds/placeholder-tiles.png", false)
}
func (b *BattleScreen) OnEnter() {
	var cols = b.tmap.Properties[property.MapColumns].(int)
	var rows = b.tmap.Properties[property.MapRows].(int)
	var tileW = b.tmap.Properties[property.MapTileWidth].(int)
	var tileH = b.tmap.Properties[property.MapTileHeight].(int)

	b.camera.X, b.camera.Y = float32(cols)/2*float32(tileW), float32(rows)/2*float32(tileH)
	b.camera.Zoom = 0.8

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
	// assets.SetTextureSmoothness("art/Battlegrounds/placeholder-tiles.png", false)

}
func (b *BattleScreen) OnUpdate() {
	var tileW = b.tmap.Properties[property.MapTileWidth].(int)
	var tileH = b.tmap.Properties[property.MapTileHeight].(int)

	b.camera.SetScreenAreaToWindow()

	if b.currentPopup == nil {
		b.camera.MouseDragAndZoomSmoothly()
		b.camera.Zoom = number.Limit(b.camera.Zoom, 0.5, 10)
	}

	// b.tmap.Draw(b.camera)
	b.camera.DrawSprites(b.tiles...)

	var x, y = b.turnManager.takingTurnUnit.Position()
	x, y = x*float32(tileW)+float32(tileW/2), y*float32(tileH)+float32(tileH/2)
	b.camera.DrawCircle(x, y, 32, palette.White)

	for _, unit := range b.ySortUnits() {
		unit.Draw(b.camera, tileW, tileH)
	}

	b.turnManager.update()

	b.hud.UpdateAndDraw(b.camera)
	if b.currentPopup != nil {
		b.currentPopup.UpdateAndDraw(b.camera)
	}

	b.handleInput()
}
func (b *BattleScreen) OnExit() {
}

func (b *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	b.units = collection.Join(teamA, teamB)
	b.spawnUnits(teamA, false, "BattleSpawnsTeamA")
	b.spawnUnits(teamB, true, "BattleSpawnsTeamB")

	var pathMapLayers = b.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")
	var pathMap *geometry.ShapeGrid
	if len(pathMapLayers) > 0 {
		pathMap = pathMapLayers[0].ExtractShapeGrid()
	}

	b.turnManager.startBattle(teamA, teamB, playerIsTeamA, pathMap)
}

//=================================================================
// private

func (b *BattleScreen) spawnUnits(units []*unit.Unit, flip bool, layerClass string) {
	var tileW = b.tmap.Properties[property.MapTileWidth].(int)
	var tileH = b.tmap.Properties[property.MapTileHeight].(int)
	var spawns = b.tmap.FindLayersBy(property.LayerClass, layerClass)[0].ExtractPoints()
	if len(units) > len(spawns) {
		return
	}
	for i, u := range units {
		var x = int(spawns[i][0] / float32(tileW))
		var y = int(spawns[i][1] / float32(tileH))
		u.Spawn(float32(x), float32(y), flip)
	}
}
func (b *BattleScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.Escape) {
		screens.Enter(global.ScreenWorld, false)
	} else if keyboard.IsKeyJustPressed(key.L) {
		b.currentPopup = global.TogglePopup(b.hud, b.currentPopup, b.loot)
	}
}
func (b *BattleScreen) ySortUnits() []*unit.Unit {
	var ySorted = make(map[float32][]*unit.Unit, len(b.units))

	for _, unit := range b.units {
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
