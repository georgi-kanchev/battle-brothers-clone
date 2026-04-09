package battle

import (
	"game/code/global"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
)

type BattleScreen struct {
	path   string // cached for reload
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	ground, above *graphics.TileMap
	pathMap       *geometry.ShapeGrid

	spawnsA, spawnsB []float32

	units       []*unit.Unit
	hoveredUnit *unit.Unit

	team1, team2  []*unit.Unit
	playerIsTeam1 bool
	order         []*unit.Unit
	states        *flow.StateMachine
	curTurnIndex  int
	curTurnIsTeam1 bool

	curMoveRangeCells [][2]int
	curMovePath       []float32
	curMoveIndex      int
}

func New(mapPath string) *BattleScreen {
	return &BattleScreen{path: mapPath, camera: graphics.NewCamera(1)}
}

//=================================================================

func (bs *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle GUI...")
	bs.hud = gui.NewFromXMLs(bs.camera, file.LoadText("data/gui/battle-hud.xml"), global.ThemesGUI)
	bs.loot = gui.NewFromXMLs(bs.camera, file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	bs.currentPopup = nil

	var sc = global.Opts.ScaleUI
	bs.hud.Scale = global.Opts.ScaleBattleHUD * sc
	bs.loot.Scale = global.Opts.ScaleBattleLoot * sc

	loading.Show("Loading:\nBattle Tiles...")
	var mapPath = "data/battlegrounds/test/map.tmx"
	var tileSetId, _ = assets.LoadTiledLayers(mapPath)

	bs.ground = graphics.NewTileMap(tileSetId, mapPath+"/Ground")
	bs.above = graphics.NewTileMap(tileSetId, mapPath+"/Above")

	bs.ground.PivotX, bs.ground.PivotY = 0, 0
	bs.above.PivotX, bs.above.PivotY = 0, 0

	var spawnsA = graphics.NewTileMap(tileSetId, mapPath+"/SpawnsA")
	spawnsA.PivotX, spawnsA.PivotY = 0, 0
	bs.spawnsA = spawnsA.Points()

	var spawnsB = graphics.NewTileMap(tileSetId, mapPath+"/SpawnsB")
	spawnsB.PivotX, spawnsB.PivotY = 0, 0
	bs.spawnsB = spawnsB.Points()

	var w, h = bs.ground.Size()
	global.BattleTileWidth, global.BattleTileHeight = bs.ground.SizeTile()
	global.BattleColumns, global.BattleRows = float32(w), float32(h)

	bs.pathMap = geometry.NewShapeGrid(int(global.BattleTileWidth), int(global.BattleTileHeight))
}
func (bs *BattleScreen) OnEnter() {
	bs.camera.X = global.BattleColumns / 2 * global.BattleTileWidth
	bs.camera.Y = global.BattleRows / 2 * global.BattleTileHeight
	bs.camera.Zoom = 0.8

	for _, id := range assets.LoadedTextureIds() { // probably shouldn't be here
		assets.SetTextureSmoothness(id, true)
	}
}
func (bs *BattleScreen) OnUpdate() {
	if bs.currentPopup == nil {
		bs.camera.MouseDragAndZoomSmoothly()
		bs.camera.Zoom = number.Limit(bs.camera.Zoom, 0.5, 10)
	}

	bs.camera.DrawTileMaps(bs.ground, bs.above)

	bs.updateUnits()

	bs.hud.UpdateAndDraw()
	if bs.currentPopup != nil {
		bs.currentPopup.UpdateAndDraw()
	}

	bs.handleInput()

	global.TryShowFPS(bs.camera)
}
func (bs *BattleScreen) OnExit() {
}

func (bs *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	bs.units = collection.Join(teamA, teamB)
	bs.team1, bs.team2 = teamA, teamB
	bs.playerIsTeam1 = playerIsTeamA
	bs.states = flow.NewStateMachine()
	bs.spawnUnits(bs.spawnsA, teamA)
	bs.spawnUnits(bs.spawnsB, teamB)
	bs.startBattle()
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
