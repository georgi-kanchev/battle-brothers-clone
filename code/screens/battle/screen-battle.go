package battle

import (
	"game/code/global"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type BattleScreen struct {
	path   string // cached for reload
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	tmap    *tiled.Map
	tiles   []*graphics.Sprite
	pathMap *geometry.ShapeGrid

	unitManager *unitManager
}

func New(mapPath string) *BattleScreen {
	return &BattleScreen{path: mapPath, camera: graphics.NewCamera(1)}
}

//=================================================================

func (bs *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	bs.unitManager = newUnitManager(teamA, teamB)
	bs.unitManager.spawnAll(bs.tmap, teamA, "BattleSpawnsTeamA")
	bs.unitManager.spawnAll(bs.tmap, teamB, "BattleSpawnsTeamB")
	bs.unitManager.turnManager.startBattle(teamA, teamB, playerIsTeamA)
}

func (bs *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle Map...")
	bs.tmap = tiled.NewMap(assets.LoadTiledMap(bs.path), global.Project)
	loading.Show("Loading:\nBattle GUI...")
	bs.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.ThemesGUI)
	bs.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	bs.currentPopup = nil

	var sc = global.Opts.ScaleUI
	bs.hud.Scale = global.Opts.ScaleBattleHUD * sc
	bs.loot.Scale = global.Opts.ScaleBattleLoot * sc

	loading.Show("Processing:\nBattle data...")
	var layers = bs.tmap.FindLayersBy(property.LayerClass, "BattleMap")
	layers = append(layers, bs.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")...)
	for _, l := range layers {
		bs.tiles = append(bs.tiles, l.ExtractSprites()...)
	}

	var atlasId = bs.tmap.Tilesets[0].Properties["atlasId"].(string)
	for i := range 31 {
		var tileId = text.New(atlasId, "/", i)
		assets.SetTextureAtlasTile(atlasId, tileId, float32(i), 1, 1, 1, 0, false)
	}
}
func (bs *BattleScreen) OnEnter() {
	global.BattleTileWidth = float32(bs.tmap.Properties[property.MapTileWidth].(int))
	global.BattleTileHeight = float32(bs.tmap.Properties[property.MapTileHeight].(int))
	global.BattleTileColumns = float32(bs.tmap.Properties[property.MapColumns].(int))
	global.BattleTileRows = float32(bs.tmap.Properties[property.MapRows].(int))

	bs.camera.X = global.BattleTileColumns / 2 * global.BattleTileWidth
	bs.camera.Y = global.BattleTileRows / 2 * global.BattleTileHeight
	bs.camera.Zoom = 0.8

	for _, id := range assets.LoadedTextureIds() { // probably shouldn't be here
		assets.SetTextureSmoothness(id, true)
	}
}
func (bs *BattleScreen) OnUpdate() {
	bs.camera.SetScreenAreaToWindow()

	if bs.currentPopup == nil {
		bs.camera.MouseDragAndZoomSmoothly()
		bs.camera.Zoom = number.Limit(bs.camera.Zoom, 0.5, 10)
	}

	bs.camera.DrawSprites(bs.tiles...)

	bs.unitManager.update()

	bs.hud.UpdateAndDraw(bs.camera)
	if bs.currentPopup != nil {
		bs.currentPopup.UpdateAndDraw(bs.camera)
	}

	bs.handleInput()

	global.TryShowFPS(bs.camera)
}

func (bs *BattleScreen) OnExit() {
}

//=================================================================
// private

func (bs *BattleScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.Escape) {
		screens.Enter(global.ScreenWorld, false)
	} else if keyboard.IsKeyJustPressed(key.L) {
		bs.currentPopup = condition.If(bs.currentPopup == bs.loot, nil, bs.loot)
	}
}
func (bs *BattleScreen) recalculatePathMap() {
	var pathMapLayers = bs.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")
	if len(pathMapLayers) > 0 {
		bs.pathMap = pathMapLayers[0].ExtractShapeGrid()
	}

	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	for _, unit := range bs.unitManager.units {
		if unit != bs.unitManager.turnManager.unitActing() {
			var ux, uy = unit.Position()
			bs.pathMap.SetAtCell(int(ux/tw), int(uy/th), geometry.NewShapeQuad(tw/2, th/2, 0.5, 0.5))
		}
	}
}
