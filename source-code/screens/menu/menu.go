package menu

import (
	"game/source-code/global"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/window"
)

type Menu struct {
	camera *graphics.Camera

	bgr, knight, logo *graphics.Sprite

	hud, play, load, options, currentPopup *gui.GUI
}

func New() *Menu {
	var menu = &Menu{camera: graphics.NewCamera(1)}
	return menu
}

//=================================================================

func (menu *Menu) OnLoad() {
	menu.hud = gui.NewFromXMLs(file.LoadText("data/gui/menu-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	// menu.play = gui.NewFromXMLs(file.LoadText("data/gui/menu-play.xml"), global.ThemesGUI)
	// menu.load = gui.NewFromXMLs(file.LoadText("data/gui/menu-load.xml"), global.ThemesGUI)
	menu.options = gui.NewFromXMLs(file.LoadText("data/gui/menu-options.xml"), global.ThemesGUI)
	menu.currentPopup = nil

	var bgr = assets.LoadTexture("art/UI/Titlescreen/background.PNG")
	var knight = assets.LoadTexture("art/UI/Titlescreen/knight.PNG")
	var logo = assets.LoadTexture("art/UI/Titlescreen/logo.PNG")
	menu.bgr = graphics.NewSprite(bgr, 0, 0)
	menu.knight = graphics.NewSprite(knight, -500, 0)
	menu.logo = graphics.NewSprite(logo, 0, 0)
	assets.SetTextureSmoothness(bgr, true)
	assets.SetTextureSmoothness(knight, true)
	assets.SetTextureSmoothness(logo, true)

	menu.logo.ScaleX, menu.logo.ScaleY = 1.5, 1.5
	menu.knight.ScaleX, menu.knight.ScaleY = 1.6, 1.6
}
func (menu *Menu) OnEnter() {
}
func (menu *Menu) OnUpdate() {
	menu.camera.SetScreenAreaToWindow()

	positionElements(menu)

	menu.hud.UpdateAndDraw(menu.camera)
	if menu.currentPopup != nil {
		menu.currentPopup.UpdateAndDraw(menu.camera)
	}

	menu.handleInput()
}

func (menu *Menu) OnExit() {
}

//=================================================================
// private

func positionElements(menu *Menu) {
	var rx, ry = menu.camera.PointFromEdge(1, 0.5)
	menu.logo.X, menu.logo.Y = rx-menu.logo.Width/2-250, ry-menu.logo.Height/2-150
	menu.bgr.CameraFill(menu.camera)
	menu.camera.DrawSprites(menu.bgr, menu.knight, menu.logo)
}
func (menu *Menu) handleInput() {
	if menu.hud.IsButtonJustClicked("new", menu.camera) {
		screens.Enter(global.ScreenWorld, false)
	} else if menu.hud.IsButtonJustClicked("options", menu.camera) {
		menu.currentPopup = global.TogglePopup(menu.hud, menu.currentPopup, menu.options)
	} else if menu.hud.IsButtonJustClicked("quit", menu.camera) {
		window.Close()
	} else if keyboard.IsKeyJustPressed(key.Escape) {
		menu.currentPopup = global.TogglePopup(menu.hud, menu.currentPopup, menu.currentPopup)
	}
}
