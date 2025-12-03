package main

import (
	"pure-game-kit/data/assets"
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

	var _, iconIds = assets.LoadDefaultAtlasIcons()

	var worldHud = gui.NewElements(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("panel", f.Color, "43 33 23 255", f.FrameColor, "86 66 46 255", f.FrameSize, "-5",
			f.TextLineHeight, "25", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5"),
		gui.Theme("resource", f.TextLineHeight, "20", f.Width, "80", f.Height, "40",
			f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5", f.TextLineGap, "-1"),
		//=================================================================
		gui.Container("resources", d.CameraLeftX+"+5", d.CameraTopY+"+5", "400", "40"),
		gui.Visual("resources-background", f.ThemeId, "panel", f.FillContainer, ""),
		gui.Visual("gold", f.ThemeId, "resource", f.Text, "^^ 90 523", f.TextEmbeddedAssetId1, iconIds[1]),
		gui.Visual("food", f.ThemeId, "resource", f.Text, "^^ 188", f.TextEmbeddedAssetId1, iconIds[2]),
		gui.Visual("tools", f.ThemeId, "resource", f.Text, "^^ 77", f.TextEmbeddedAssetId1, iconIds[3]),
		gui.Visual("ammo", f.ThemeId, "resource", f.Text, "^^ 68", f.TextEmbeddedAssetId1, iconIds[4]),
		gui.Visual("medicine", f.ThemeId, "resource", f.Text, "^^ 29", f.TextEmbeddedAssetId1, iconIds[5]),
		//=================================================================
		gui.Container("goal", d.TargetLeftX, d.TargetBottomY, d.TargetWidth, "40", f.TargetId, "resources"),
		gui.Visual("goal-background", f.ThemeId, "panel", f.FillContainer, "", f.Text, "goal"),
		//=================================================================
		gui.Container("time", d.CameraCenterX+"-50", d.CameraTopY+"+5", "80", "80"),
		gui.Visual("time-background", f.ThemeId, "panel", f.FillContainer, "", f.Text, "day 16",
			f.TextAlignmentY, "0"),
		//=================================================================
		gui.Container("menus", d.CameraRightX+"-405", d.CameraTopY+"+5", "400", "40"),
		gui.Visual("menus-background", f.ThemeId, "panel", f.FillContainer, ""),
		//=================================================================
		gui.Container("quest", d.TargetLeftX, d.TargetBottomY, d.TargetWidth, "40", f.TargetId, "menus"),
		gui.Visual("quest-background", f.ThemeId, "panel", f.FillContainer, "", f.Text, "quest"),
	)

	worldHud.Scale = 2

	assets.LoadDefaultFont()

	var x, y, moveTargetX, moveTargetY float32

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
		camera.Zoom = number.Limit(camera.Zoom, 0.1, 8)

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

		camera.DrawCircle(x, y, 20, color.Black)
		camera.DrawCircle(x, y, 16, color.Cyan)

		if !worldHud.IsAnyHovered(camera) {
			if mouse.IsButtonPressed(button.Left) {
				worldHud.IsAnyHovered(camera)
				moveTargetX, moveTargetY = camera.MousePosition()
			}
			if mouse.IsButtonJustPressed(button.Right) {
				x, y = camera.MousePosition()
				moveTargetX, moveTargetY = x, y
			}
		}

		worldHud.UpdateAndDraw(camera)
	}
}
