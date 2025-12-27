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
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

type Menu struct {
	camera *graphics.Camera

	bgr, logo *graphics.Sprite

	hud, play, load, options, currentPopup *gui.GUI
}

func New() *Menu {
	var menu = &Menu{camera: graphics.NewCamera(1)}
	return menu
}

//=================================================================

func (menu *Menu) OnLoad() {
	menu.hud = gui.NewFromXMLs(file.LoadText("data/gui/menu-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	menu.options = gui.NewFromXMLs(file.LoadText("data/gui/menu-options.xml"), global.ThemesGUI)
	menu.currentPopup = nil

	var bgr = assets.LoadTexture("art/UI/Titlescreen/bgr.png")
	var logo = assets.LoadTexture("art/UI/Titlescreen/logo.PNG")
	menu.bgr = graphics.NewSprite(bgr, 0, 0)
	menu.logo = graphics.NewSprite(logo, 0, 0)
	assets.SetTextureSmoothness(bgr, true)
	assets.SetTextureSmoothness(logo, true)

	menu.logo.ScaleX, menu.logo.ScaleY = 0.8, 0.8
	menu.logo.PivotX, menu.logo.PivotY = 1, 1

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (menu *Menu) OnEnter() {
}
func (menu *Menu) OnUpdate() {
	menu.camera.DrawColor(color.RGB(8, 3, 4))
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
	var sc = menu.bgr.ScaleX

	menu.bgr.CameraFit(menu.camera)
	menu.bgr.X -= 100 * sc
	menu.bgr.ScaleX, menu.bgr.ScaleY = menu.bgr.ScaleX*1.1, menu.bgr.ScaleY*1.1
	menu.logo.X, menu.logo.Y = rx-80*sc, ry
	menu.logo.ScaleX, menu.logo.ScaleY = sc, sc
	menu.hud.Scale = sc * 1.5
	menu.camera.DrawSprites(menu.bgr, menu.logo)
}
func (menu *Menu) handleInput() {
	if menu.hud.IsButtonJustClicked("new", menu.camera) {
		screens.Enter(global.ScreenWorld, false)
	} else if menu.hud.IsButtonJustClicked("options", menu.camera) {
		menu.currentPopup = global.TogglePopup(menu.hud, menu.currentPopup, menu.options)
	} else if menu.hud.IsButtonJustClicked("quit", menu.camera) {
		window.Close()
	} else if menu.currentPopup != nil && keyboard.IsKeyJustPressed(key.Escape) {
		menu.currentPopup = global.TogglePopup(menu.hud, menu.currentPopup, menu.currentPopup)
	}
}
