package battle

import (
	"game/source-code/global"
	"game/source-code/screens/loading"
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
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type BattleScreen struct {
	path   string
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	tmap    *tiled.Map
	tiles   []*graphics.Sprite
	pathMap *geometry.ShapeGrid

	units []*unit.Unit

	turnManager *turnManager

	hoveredWalkRange bool
	hoveredCell      [2]int
	hoveredUnit      *unit.Unit
}

func New(mapPath string) *BattleScreen {
	var battle = &BattleScreen{path: mapPath, camera: graphics.NewCamera(1), turnManager: newTurnManager()}
	return battle
}

//=================================================================

func (b *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle Map...")
	b.tmap = tiled.NewMap(assets.LoadTiledMap(b.path), global.Project)
	loading.Show("Loading:\nBattle GUI...")
	b.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	b.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	b.currentPopup = nil

	var sc = global.Options.ScaleUI.Master
	b.hud.Scale = global.Options.ScaleUI.Battle.HUD * sc
	b.loot.Scale = global.Options.ScaleUI.Battle.Loot * sc

	loading.Show("Processing:\nBattle data...")
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
	var tileW, tileH = b.tileSize()

	b.camera.SetScreenAreaToWindow()

	if b.currentPopup == nil {
		b.camera.MouseDragAndZoomSmoothly()
		b.camera.Zoom = number.Limit(b.camera.Zoom, 0.5, 10)
	}

	// b.tmap.Draw(b.camera)
	b.camera.DrawSprites(b.tiles...)

	b.calculateHoverInfo()

	var ySortedUnits = b.ySortUnits()
	b.turnManager.states.UpdateCurrentState()

	b.drawBehindUnits()

	for _, unit := range ySortedUnits {
		unit.Draw(b.camera, tileW, tileH)
	}
	b.drawAboveUnits()

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

	b.turnManager.startBattle(teamA, teamB, playerIsTeamA)
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
func (b *BattleScreen) calculateHoverInfo() {
	var battle = screens.Current().(*BattleScreen)
	var tileW, tileH = battle.tileSize()
	var mx, my = battle.camera.MousePosition()

	b.hoveredCell[0], b.hoveredCell[1] = int(mx/tileW), int(my/tileH)

	b.hoveredUnit = nil
	for _, unit := range b.units {
		if unit.IsHovered(b.camera, mx/tileW, my/tileH) {
			b.hoveredUnit = unit
			break
		}
	}
}
func (b *BattleScreen) recalculatePathMap() {
	var pathMapLayers = b.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")
	if len(pathMapLayers) > 0 {
		b.pathMap = pathMapLayers[0].ExtractShapeGrid()
	}

	var tileW, tileH = b.tileSize()
	for _, unit := range b.units {
		if unit == b.turnManager.unit() {
			continue
		}
		var ux, uy = unit.Position(tileW, tileH)
		b.pathMap.SetAtCell(int(ux/tileW), int(uy/tileH), geometry.NewShapeRectangle(tileW/2, tileH/2, 0.5, 0.5))
	}
}

func (b *BattleScreen) drawBehindUnits() {
	var tileW, tileH = b.tileSize()

	for _, cell := range b.turnManager.curWalkRangeCells {
		var cx, cy = float32(cell[0]) * tileW, float32(cell[1]) * tileH
		b.camera.DrawQuad(cx, cy, float32(tileW), float32(tileH), 0, color.FadeOut(palette.Red, 0.5))
	}

	if b.hoveredUnit != nil {
		var ux, uy = b.hoveredUnit.Position(tileW, tileH)
		b.camera.DrawQuad(ux-tileW/2, uy-tileH/2, 64, 64, 0, color.FadeOut(palette.White, 0.75))
	}

	var ux, uy = b.turnManager.unit().Position(tileW, tileH)
	b.camera.DrawQuadFrame(ux-tileW/2, uy-tileH/2, tileW, tileH, 0, -2, palette.Azure)

	var hx, hy = float32(b.hoveredCell[0]) * tileW, float32(b.hoveredCell[1]) * tileH
	b.camera.DrawQuadFrame(hx, hy, tileW, tileH, 0, -1, palette.White)
}
func (b *BattleScreen) drawUnitStats(description string, unit *unit.Unit) {
	var lineHeight = 80 / b.camera.Zoom
	var txt = text.New(
		description, "\n",
		"Initiative: ", unit.Initiative, "\n",
		"Movement: ", unit.Movement, "\n",
	)
	var lines = len(text.Split(txt, "\n"))
	var blx, bly = b.camera.PointFromEdge(0, 1)
	b.camera.DrawText("", txt, blx+50/b.camera.Zoom, bly-lineHeight*float32(lines), lineHeight, palette.White)
}
func (b *BattleScreen) drawAboveUnits() {
	var tileW, tileH = b.tileSize()
	var unit = b.turnManager.order[b.turnManager.curIndex]
	var mx, my = b.camera.MousePosition()
	var ux, uy = unit.Position(tileW, tileH)
	var path = b.pathMap.FindPathDiagonally(ux, uy, mx, my, false)

	if len(path) > 0 {
		b.camera.DrawLinesPath(6, palette.Azure, path...)
		b.camera.DrawPoints(3, palette.Red, path...)

		var x, y = path[len(path)-1][0], path[len(path)-1][1]
		var txt = text.New(unit.CalculateMovementPoints(path, tileW, tileH), "/", unit.Movement)
		b.camera.DrawText("", txt, x, y, 50/b.camera.Zoom, palette.White)
	}

	if b.hoveredUnit != nil {
		b.drawUnitStats("Hovered Unit", b.hoveredUnit)
	} else {
		b.drawUnitStats("Unit taking turn", b.turnManager.unit())
	}
}

func (b *BattleScreen) ySortUnits() []*unit.Unit {
	var ySorted = make(map[float32][]*unit.Unit, len(b.units))

	for _, unit := range b.units {
		var _, y = unit.Cell()
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
func (b *BattleScreen) tileSize() (width, height float32) {
	var tileW = b.tmap.Properties[property.MapTileWidth].(int)
	var tileH = b.tmap.Properties[property.MapTileHeight].(int)
	return float32(tileW), float32(tileH)
}
