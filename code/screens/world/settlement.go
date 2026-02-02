package world

import (
	"game/code/global"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (ws *WorldScreen) handleSettlementPopup() {
	var name = ws.playerParty.goingToSettlement.Properties[property.ObjectName].(string)
	ws.settlement.SetField("title-label", field.Text, name)

	if ws.settlement.IsButtonJustClicked("exit-btn", ws.camera) {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.settlement)
		ws.playerParty.goingToSettlement = nil
	}

	if ws.settlement.IsButtonJustClicked("rest", ws.camera) {
		ws.playerParty.isResting = true
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.settlement)
	} else if ws.settlement.IsButtonJustClicked("market", ws.camera) {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.market)
	}
}
