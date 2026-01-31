package world

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
)

func (ws *WorldScreen) handleResting() {
	ws.hud.SetField("rest", field.Hidden, condition.If(ws.playerParty.isResting, "", "1"))

	if !ws.playerParty.isResting {
		return
	}

	var moveCancel = mouse.IsAnyButtonJustPressed() && !ws.hud.IsAnyHovered(ws.camera)

	if moveCancel {
		ws.playerParty.isResting = false
		ws.playerParty.goingToSettlement = nil
	}
}
