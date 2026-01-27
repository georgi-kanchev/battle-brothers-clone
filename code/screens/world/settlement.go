package world

import (
	"game/code/global"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (ws *WorldScreen) handleSettlementPopup() {
	var name = ws.playerParty.goingToSettlement.Properties[property.ObjectName].(string)
	ws.settlement.SetField("settlement-title-label", field.Text, name)

	if ws.settlement.IsButtonJustClicked("settlement-exit-btn", ws.camera) {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.settlement)
		ws.playerParty.goingToSettlement = nil
	}
}
