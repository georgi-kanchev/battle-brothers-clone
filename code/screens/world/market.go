package world

func (ws *WorldScreen) handleMarket() {
	if ws.settlement.IsButtonJustClicked("exit-btn", ws.camera) {
		ws.currentPopup = ws.settlement
	}
}
