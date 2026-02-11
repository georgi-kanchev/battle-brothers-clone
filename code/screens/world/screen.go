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
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
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
	settlements  *tiled.Layer

	tmap      *tiled.Map
	mapLayers []*tiled.Layer
	solids    []*geometry.Shape
	roads     [][2]float32
}

func New(path string) *WorldScreen {
	return &WorldScreen{path: path, camera: graphics.NewCamera(1), time: 60 * 3}
}

//=================================================================

func (ws *WorldScreen) OnLoad() {
	loading.Show("Loading:\nWorld Map...")
	ws.tmap = tiled.NewMap(assets.LoadTiledMap(ws.path), global.Project)
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

	ws.hud = gui.NewFromXMLs(hud, dim, themes)
	ws.events = gui.NewFromXMLs(events, themes)
	ws.inventory = gui.NewFromXMLs(dim, wide, inventory, x, themes)
	ws.settlement = gui.NewFromXMLs(dim, narrow, settlement, title, x, themes)
	ws.market = gui.NewFromXMLs(dim, wide, market, title, x, themes)
	ws.favors = gui.NewFromXMLs(dim, narrow, favors, title, x, themes)
	ws.recruit = gui.NewFromXMLs(dim, narrow, recruit, title, x, themes)
	ws.tavern = gui.NewFromXMLs(dim, narrow, tavern, title, x, themes)
	ws.currentPopup = nil

	var sc = global.Options.ScaleUI
	ws.hud.Scale = global.Options.ScaleWorldHUD * sc
	ws.inventory.Scale = global.Options.ScaleWorldInventory * sc
	ws.events.Scale = global.Options.ScaleWorldEvents * sc
	ws.settlement.Scale = global.Options.ScaleWorldSettlement * sc
	ws.market.Scale = global.Options.ScaleWorldSettlementMarket * sc
	ws.favors.Scale = global.Options.ScaleWorldSettlementMarket * sc
	ws.recruit.Scale = global.Options.ScaleWorldSettlementRecruit * sc
	ws.tavern.Scale = global.Options.ScaleWorldSettlementTavern * sc

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

	unit.Names = text.Split(file.LoadText("data/names.txt"), " ")
	unit.Nicknames = text.Split(file.LoadText("data/nicknames.txt"), " ")

	loading.Show("Processing:\nWorld data...")
	var solidLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldSolids")
	var roadLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldRoads")
	var settlements = ws.tmap.FindLayersBy(property.LayerClass, "WorldSettlements")
	ws.mapLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldMap")
	for _, s := range solidLayers {
		ws.solids = append(ws.solids, s.ExtractShapes()...)
	}
	for _, r := range roadLayers {
		ws.roads = append(ws.roads, r.ExtractLines()...)
	}
	if len(settlements) > 0 {
		ws.settlements = settlements[0]
	}

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (ws *WorldScreen) OnEnter() {
	var units = []*unit.Unit{unit.New(), unit.New(), unit.New(), unit.New(), unit.New(), unit.New()}
	ws.playerParty = NewParty(units, 2250, 1530, true)
}
func (ws *WorldScreen) OnUpdate() {
	ws.camera.SetScreenAreaToWindow()

	//world.tmap.Draw(world.camera)
	for _, m := range ws.mapLayers {
		m.Draw(ws.camera)
	}

	ws.handleResting()
	ws.playerParty.Update()
	for _, party := range ws.otherParties {
		party.Update()
	}
	ws.handleDayNightCycle()
	ws.handleInput()

	ws.hud.UpdateAndDraw(ws.camera)
	ws.hud.SetField("popup-dim-bgr", field.Hidden, condition.If(ws.currentPopup == nil, "1", ""))

	if ws.currentPopup != nil {
		ws.currentPopup.UpdateAndDraw(ws.camera)
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
}

func (ws *WorldScreen) OnExit() {
}

//=================================================================
// private

var teamA = []*unit.Unit{}
var teamB = []*unit.Unit{}

func (ws *WorldScreen) handleInput() {
	if ws.currentPopup != nil {
		return
	}

	if ws.hud.IsButtonJustClicked("inventory", ws.camera) {
		ws.currentPopup = ws.inventory
	}

	if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
		var scr = screens.Current().(*battle.BattleScreen)
		scr.Prepare(teamA, teamB, true)
	} else if keyboard.IsKeyJustPressed(key.E) {
		ws.currentPopup = ws.events
	}

	if ws.hud.IsButtonJustClicked("main-menu", ws.camera) {
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
	if from.IsButtonJustClicked("exit-btn", ws.camera) ||
		from.IsButtonJustClicked("popup-dim-bgr", ws.camera) {
		ws.currentPopup = to

		if andDo != nil {
			andDo()
		}
	}
}
