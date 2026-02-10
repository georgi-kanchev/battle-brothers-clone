package world

import "pure-game-kit/gui/field"

func (ws *WorldScreen) handleTavernPopup() {
	ws.tavern.SetField("title-bgr", field.Text, "Tavern")
	ws.tryExitPopup(ws.tavern, ws.settlement, nil)
}
