package world

import (
	"game/code/global"
	"game/code/screens/battle"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/data/folder"
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/text"
)

type WorldScreen struct {
	path   string
	camera *graphics.Camera

	hud, inventory, events, currentPopup        *gui.GUI
	settlement, market, favors, recruit, tavern *gui.GUI

	resultingCursorNonGUI int

	time       float32
	timeCircle *graphics.Sprite

	playerParty  *Party
	otherParties []*Party

	chunks              []*graphics.Sprite
	solids, settlements []*geometry.Shape
	roads               []float32
}

func New(path string) *WorldScreen {
	return &WorldScreen{path: path, camera: graphics.NewCamera(1), time: 60 * 3}
}

//=================================================================

func (ws *WorldScreen) OnLoad() {
	loading.Show("Loading:\nWorld Map...")
	//ws.scene = tiled.NewScene(assets.LoadTiledMap(ws.path), global.Project)
	loading.Show("Loading:\nWorld GUI...")

	var narrow, wide = global.PopupNarrowGUI, global.PopupWideGUI
	var dim, x, title, themes = global.DimGUI, global.XBtnGUI, global.TitleGUI, global.ThemesGUI
	var hud = file.LoadText("data/gui/world-hud.xml")
	var events = file.LoadText("data/gui/world-events.xml")
	var inventory = file.LoadText("data/gui/world-inventory.xml")
	var settlement = file.LoadText("data/gui/world-settlement.xml")
	var market = file.LoadText("data/gui/world-settlement-market.xml")
	var favors = file.LoadText("data/gui/world-settlement-favors.xml")
	var recruit = file.LoadText("data/gui/world-settlement-recruit.xml")
	var tavern = file.LoadText("data/gui/world-settlement-tavern.xml")

	ws.hud = gui.NewFromXMLs(ws.camera, hud, dim, themes)
	ws.events = gui.NewFromXMLs(ws.camera, events, themes)
	ws.inventory = gui.NewFromXMLs(ws.camera, dim, wide, inventory, x, themes)
	ws.settlement = gui.NewFromXMLs(ws.camera, dim, narrow, settlement, title, x, themes)
	ws.market = gui.NewFromXMLs(ws.camera, dim, wide, market, title, x, themes)
	ws.favors = gui.NewFromXMLs(ws.camera, dim, narrow, favors, title, x, themes)
	ws.recruit = gui.NewFromXMLs(ws.camera, dim, narrow, recruit, title, x, themes)
	ws.tavern = gui.NewFromXMLs(ws.camera, dim, narrow, tavern, title, x, themes)
	ws.currentPopup = nil

	var sc = global.Opts.ScaleUI
	ws.hud.Scale = global.Opts.ScaleWorldHUD * sc
	ws.inventory.Scale = global.Opts.ScaleWorldInventory * sc
	ws.events.Scale = global.Opts.ScaleWorldEvents * sc
	ws.settlement.Scale = global.Opts.ScaleWorldSettlement * sc
	ws.market.Scale = global.Opts.ScaleWorldSettlementMarket * sc
	ws.favors.Scale = global.Opts.ScaleWorldSettlementMarket * sc
	ws.recruit.Scale = global.Opts.ScaleWorldSettlementRecruits * sc
	ws.tavern.Scale = global.Opts.ScaleWorldSettlementTavern * sc

	loading.Show("Loading:\nWorld images...")
	var timeCircle = assets.LoadTexture("art/UI/Time/time_circle.PNG")
	ws.timeCircle = graphics.NewSprite(timeCircle, 0, 0)

	assets.LoadTexture("art/UI/Time/time_top.PNG")
	assets.LoadTexture("art/UI/Buttons/btn.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_pause.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_play.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx2.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx3.PNG")

	var allAssets []string
	allAssets = append(allAssets, folder.Content("art/Character", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/hair", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/body_armor", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/head_armor", true)...)
	for _, filePath := range allAssets {
		if path.IsFile(filePath) {
			assets.LoadTexture(filePath)
		}
	}

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}

	unit.Names = text.Split(file.LoadText("data/names.txt"), " ")
	unit.Nicknames = text.Split(file.LoadText("data/nicknames.txt"), " ")

	loading.Show("Processing:\nWorld Map...")
	var chunks = folder.Content("data/worlds/01", true)
	for _, file := range chunks {
		if path.Extension(file) != ".jpg" {
			continue
		}
		var name = path.RemoveExtension(path.LastPart(file))
		var split = text.Split(name, " ")
		var assetId = assets.LoadTexture(file)
		var w, h = assets.Size(assetId)
		var col, row = text.ToNumber[int](split[0]), text.ToNumber[int](split[1])
		var spr = graphics.NewSprite(assetId, float32(col*w), float32(row*h))
		spr.Width++
		spr.Height++
		spr.PivotX, spr.PivotY = 0, 0
		ws.chunks = append(ws.chunks, spr)
	}

	loading.Show("Processing:\nWorld Map Shapes...")
	ws.roads = assets.LoadTiledPoints("data/worlds/01/points.tmx", "Roads")
	ws.settlements = geometry.NewShapes(assets.LoadTiledPoints("data/worlds/01/points.tmx", "Settlements")...)
	ws.solids = geometry.NewShapes(assets.LoadTiledPoints("data/worlds/01/points.tmx", "Solids")...)
}
func (ws *WorldScreen) OnEnter() {
	var units = []*unit.Unit{unit.New(), unit.New(), unit.New(), unit.New(), unit.New(), unit.New()}
	ws.playerParty = NewParty(units, 1900, 1300, true)
	assets.SetTextureSmoothness("art/UI/Time/time_circle.PNG", true) // force, despite setting
}
func (ws *WorldScreen) OnUpdate() {
	ws.camera.DrawSprites(ws.chunks...)

	for _, s := range ws.solids {
		ws.camera.DrawShapes(palette.Red, s.CornerPoints()...)
	}

	ws.handleResting()
	ws.playerParty.Update()
	for _, party := range ws.otherParties {
		party.Update()
	}
	ws.handleDayNightCycle()
	ws.handleInput()

	ws.hud.UpdateAndDraw()
	ws.hud.SetField("popup-dim", field.Hidden, condition.If(ws.currentPopup == nil, "1", ""))

	if ws.currentPopup != nil {
		ws.currentPopup.UpdateAndDraw()
	} else if ws.resultingCursorNonGUI != -1 {
		mouse.SetCursor(ws.resultingCursorNonGUI)
	}

	switch ws.currentPopup {
	case ws.inventory:
		ws.handleInventoryPopup()
	case ws.events:
		ws.handleEventsPopup()
	case ws.settlement:
		ws.handleSettlementPopup()
	case ws.market:
		ws.handleMarketPopup()
	case ws.favors:
		ws.handleFavorsPopup()
	case ws.recruit:
		ws.handleRecruitPopup()
	case ws.tavern:
		ws.handleTavernPopup()
	}

	global.TryShowFPS(ws.camera)
}

func (ws *WorldScreen) OnExit() {
}

//=================================================================
// private

func (ws *WorldScreen) handleInput() {
	if ws.currentPopup != nil {
		return
	}

	if ws.hud.IsButtonJustClicked("inventory") {
		ws.currentPopup = ws.inventory
	}

	if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
		var scr = screens.Current().(*battle.BattleScreen)
		var teamA = []*unit.Unit{unit.New(), unit.New(), unit.New()}
		var teamB = []*unit.Unit{unit.New(), unit.New()}
		scr.Prepare(teamA, teamB, true)
	} else if keyboard.IsKeyJustPressed(key.E) {
		ws.currentPopup = ws.events
	}

	if ws.hud.IsButtonJustClicked("main-menu") {
		var escape = keyboard.IsKeyJustPressed(key.Escape)
		var resting = ws.playerParty.isResting

		if (resting && !escape) || !resting {
			screens.Enter(global.ScreenMainMenu, false)
		} else if resting && escape {
			ws.stopResting(true)
		}
	}
}

func (ws *WorldScreen) tryExitPopup(from *gui.GUI, to *gui.GUI, andDo func()) {
	if from.IsButtonJustClicked("exit-btn") || from.IsButtonJustClicked("popup-dim-bgr") {
		ws.currentPopup = to

		if andDo != nil {
			andDo()
		}
	}
}
