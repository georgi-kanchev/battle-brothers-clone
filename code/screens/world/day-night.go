package world

import (
	"game/code/global"
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

func (w *WorldScreen) handleDayNightCycle() {
	var dayNightCycleDuration = time.FromMinutes(float32(len(clock)))
	var topX, topY = w.camera.PointFromEdge(0.5, 0)
	w.timeCircle.X, w.timeCircle.Y = topX, topY+(185*w.hud.Scale/w.camera.Zoom)
	w.timeCircle.ScaleX = -1 / w.camera.Zoom * 0.5 * w.hud.Scale
	w.timeCircle.ScaleY = 1 / w.camera.Zoom * 0.5 * w.hud.Scale

	var scrX, scrY = w.camera.PointToScreen(topX, topY)
	var scrW, scrH = 100 * w.hud.Scale, 150 * w.hud.Scale
	w.camera.Mask(int(float32(scrX)-50*w.hud.Scale), int(scrY), int(scrW), int(scrH))
	w.timeCircle.Angle = number.Map(w.time, 0, dayNightCycleDuration, 0, 360)

	w.camera.DrawSprites(w.timeCircle)
	w.camera.SetScreenAreaToWindow()

	var x0, x1 = w.hud.Field("x0", field.Value), w.hud.Field("x1", field.Value)
	var x2, x3 = w.hud.Field("x2", field.Value), w.hud.Field("x3", field.Value)
	if w.currentPopup != nil || x0 != "" {
		global.TimeScale = 0
	} else if x1 != "" {
		global.TimeScale = 1
	} else if x2 != "" {
		global.TimeScale = 2
	} else if x3 != "" {
		global.TimeScale = 4
	}

	w.time += time.FrameDelta() * global.TimeScale
	w.time = number.Wrap(w.time, 0, dayNightCycleDuration)
	var timeOfDay = number.Round(w.time / 60)
	var timeOfDayIndex = int(timeOfDay) % len(clock)
	w.hud.SetField("time-word", field.Text, clock[timeOfDayIndex])

	var curColor, nextColor = overlayColors[timeOfDayIndex], overlayColors[int(timeOfDay+1)%len(clock)]
	var progress = number.Map(float32(timeOfDayIndex)-w.time/60, -0.5, 0.5, 1, 0)
	var col = color.Fade(curColor, nextColor, progress)

	w.camera.DrawColor(col)
}
