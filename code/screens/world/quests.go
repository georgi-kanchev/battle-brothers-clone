package world

func (ws *WorldScreen) handleQuestsPopup() {
	if ws.quests.IsButtonJustClicked("exit-btn", ws.camera) ||
		ws.quests.IsButtonJustClicked("popup-dim-bgr", ws.camera) {
		ws.currentPopup = ws.settlement
	}
}
