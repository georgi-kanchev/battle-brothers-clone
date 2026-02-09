package world

import (
	"game/code/global"
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
)

var overlayColors = []uint{
	color.RGBA(0, 10, 30, 160),   // midnight
	color.RGBA(0, 10, 30, 160),   // night
	color.RGBA(255, 120, 60, 51), // dawn
	0, 0, 0,                      // morning, noon, afternoon
	color.RGBA(255, 100, 40, 46), // dusk
	color.RGBA(0, 10, 30, 160),   // evening
}
var clock = []string{"Midnight", "Night", "Dawn", "Morning", "Noon", "Afternoon", "Dusk", "Evening"}

func (ws *WorldScreen) handleDayNightCycle() {
	var dayNightCycleDuration = time.FromMinutes(float32(len(clock)))
	var topX, topY = ws.camera.PointFromEdge(0.5, 0)
	ws.timeCircle.X, ws.timeCircle.Y = topX, topY+(185*ws.hud.Scale/ws.camera.Zoom)
	ws.timeCircle.ScaleX = -1 / ws.camera.Zoom * 0.5 * ws.hud.Scale
	ws.timeCircle.ScaleY = 1 / ws.camera.Zoom * 0.5 * ws.hud.Scale

	var scrX, scrY = ws.camera.PointToScreen(topX, topY)
	var scrW, scrH = 100 * ws.hud.Scale, 150 * ws.hud.Scale
	ws.camera.Mask(int(float32(scrX)-50*ws.hud.Scale), int(scrY), int(scrW), int(scrH))
	ws.timeCircle.Angle = number.Map(ws.time, 0, dayNightCycleDuration, 0, 360)

	ws.camera.DrawSprites(ws.timeCircle)
	ws.camera.SetScreenAreaToWindow()

	var x0, x1 = ws.hud.Field("x0", field.Value, ws.camera), ws.hud.Field("x1", field.Value, ws.camera)
	var x2, x3 = ws.hud.Field("x2", field.Value, ws.camera), ws.hud.Field("x3", field.Value, ws.camera)
	if ws.currentPopup != nil || x0 != "" {
		global.TimeScale = 0
	} else if x1 != "" {
		global.TimeScale = 1
	} else if x2 != "" {
		global.TimeScale = 2
	} else if x3 != "" {
		global.TimeScale = 4
	}

	ws.time += time.FrameDelta() * global.TimeScale
	ws.time = number.Wrap(ws.time, 0, dayNightCycleDuration)
	var timeOfDay = number.Round(ws.time / 60)
	var timeOfDayIndex = int(timeOfDay) % len(clock)
	ws.hud.SetField("time-word", field.Text, clock[timeOfDayIndex])
	ws.hud.SetField("paused", field.Hidden, condition.If(global.TimeScale == 0, "", "1"))

	var curColor, nextColor = overlayColors[timeOfDayIndex], overlayColors[int(timeOfDay+1)%len(clock)]
	var progress = number.Map(float32(timeOfDayIndex)-ws.time/60, -0.5, 0.5, 1, 0)
	var col = color.Fade(curColor, nextColor, progress)

	ws.camera.DrawColor(col)
}
