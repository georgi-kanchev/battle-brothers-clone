package world

func (ws *WorldScreen) handleQuestsPopup() {
	// cool effect idea - pressing the + should make the whole favor slide to the left
	// (towards the notebook) before updating all favors in the popup

	if ws.quests.IsButtonJustClicked("exit-btn", ws.camera) ||
		ws.quests.IsButtonJustClicked("popup-dim-bgr", ws.camera) {
		ws.currentPopup = ws.settlement
	}
}
