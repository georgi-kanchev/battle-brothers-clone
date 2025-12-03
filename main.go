package main

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	d "pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func main() {
	var camera = graphics.NewCamera(1)
	var assetId = assets.LoadTiledMap("data/worlds/test/map.tmx")
	var projectId = assets.LoadTiledProject("project.tiled-project")
	var project = tiled.NewProject(projectId)
	var mapp = tiled.NewMap(assetId, project)
	var worldHud = initWorldHud()

	assets.LoadDefaultAtlasIcons()
	assets.LoadDefaultFont()
	assets.LoadDefaultTexture()

	var x, y, moveTargetX, moveTargetY float32

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		
		//=================================================================
		// map

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
		}

		mapp.Draw(camera)

		//=================================================================
		// player

		x, y = point.MoveToPoint(x, y, moveTargetX, moveTargetY, 50*time.FrameDelta())
		camera.X, camera.Y = x, y

		camera.DrawTexture("", x-15, y-15, 30, 30, 0, color.Cyan)

		if keyboard.IsKeyJustPressed(key.I) {
			var hidden = condition.If(worldHud.Field("inventory", f.Hidden) == "", "1", "")
			worldHud.SetField("popup-dim", f.Hidden, hidden)
			worldHud.SetField("inventory", f.Hidden, hidden)
			time.SetScale(condition.If(hidden == "1", float32(1), 0))
		}

		if !worldHud.IsAnyHovered(camera) {
			camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
			camera.Zoom = number.Limit(camera.Zoom, 0.1, 8)

			if mouse.IsButtonPressed(button.Left) {
				worldHud.IsAnyHovered(camera)
				moveTargetX, moveTargetY = camera.MousePosition()
			}
			if mouse.IsButtonJustPressed(button.Right) {
				x, y = camera.MousePosition()
				moveTargetX, moveTargetY = x, y
			}
		}

		//=================================================================
		// gui

		worldHud.UpdateAndDraw(camera)
	}
}

func initWorldHud() *gui.GUI {
	var result = gui.NewFromXML(gui.NewElementsXML(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("panel", f.Color, "43 33 23 255", f.FrameColor, "86 66 46 255", f.FrameSize, "-5",
			f.TextLineHeight, "25", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5"),
		gui.Theme("resource", f.TextLineHeight, "20", f.Width, "80", f.Height, "40",
			f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5", f.TextLineGap, "-1"),

		//=================================================================
		// top-left
		gui.Container("resources", d.CameraLeftX+"+5", d.CameraTopY+"+5", "400", "40"),
		gui.Visual("resources-bgr", f.ThemeId, "panel", f.FillContainer, ""),
		gui.Visual("gold", f.ThemeId, "resource", f.Text, "^^ 90 523", f.TextEmbeddedAssetId1, "@coins3"),
		gui.Visual("food", f.ThemeId, "resource", f.Text, "^^ 188", f.TextEmbeddedAssetId1, "@apple"),
		gui.Visual("tools", f.ThemeId, "resource", f.Text, "^^ 77", f.TextEmbeddedAssetId1, "@ingot"),
		gui.Visual("ammo", f.ThemeId, "resource", f.Text, "^^ 68", f.TextEmbeddedAssetId1, "@bow"),
		gui.Visual("medicine", f.ThemeId, "resource", f.Text, "^^ 29", f.TextEmbeddedAssetId1, "@heart"),
		//=================================================================
		gui.Container("goal", d.TargetLeftX, d.TargetBottomY, d.TargetWidth, "40", f.TargetId, "resources"),
		gui.Visual("goal-bgr", f.ThemeId, "panel", f.FillContainer, "", f.Text, "goal"),

		//=================================================================
		// top
		gui.Container("time", d.CameraCenterX+"-50", d.CameraTopY+"+5", "80", "80"),
		gui.Visual("time-bgr", f.ThemeId, "panel", f.FillContainer, "", f.Text, "day 16", f.TextAlignmentY, "0"),

		//=================================================================
		// top-right
		gui.Container("menus", d.CameraRightX+"-405", d.CameraTopY+"+5", "400", "40"),
		gui.Visual("menus-bgr", f.ThemeId, "panel", f.FillContainer, ""),
		//=================================================================
		gui.Container("quest", d.TargetLeftX, d.TargetBottomY, d.TargetWidth, "40", f.TargetId, "menus"),
		gui.Visual("quest-bgr", f.ThemeId, "panel", f.FillContainer, "", f.Text, "quest"),

		//=================================================================
		// popups
		gui.Container("popup-dim", d.CameraLeftX, d.CameraTopY, d.CameraWidth, d.CameraHeight, f.Hidden, "1"),
		gui.Visual("popup-dim-bgr", f.FillContainer, "", f.Color, "0 0 0 150"),

		//=================================================================
		// inventory
		gui.Container("inventory", d.CameraCenterX+"-500", d.CameraCenterY+"-320", "1000", "700", f.Hidden, "1"),
		//=================================================================
		gui.Container("equipment", d.TargetLeftX, d.TargetTopY, "280", "450", f.TargetId, "inventory",
			f.Hidden, d.TargetHidden),
		gui.Visual("equipment-bgr", f.ThemeId, "panel", f.FillContainer, ""),
		//=================================================================
		gui.Container("stash", d.TargetLeftX+"+280", d.TargetTopY+"+100", "440", "350", f.TargetId, "inventory",
			f.Hidden, d.TargetHidden),
		gui.Visual("stash-bgr", f.ThemeId, "panel", f.FillContainer, ""),
		//=================================================================
		gui.Container("perks", d.TargetRightX+"-280", d.TargetTopY, "280", "450", f.TargetId, "inventory",
			f.Hidden, d.TargetHidden),
		gui.Visual("perks-bgr", f.ThemeId, "panel", f.FillContainer, ""),
		//=================================================================
		gui.Container("units", d.TargetLeftX, d.TargetBottomY+"-250", "1000", "250", f.TargetId, "inventory",
			f.Hidden, d.TargetHidden),
		gui.Visual("units-bgr", f.ThemeId, "panel", f.FillContainer, ""),
	))
	result.Scale = 2
	return result
}
