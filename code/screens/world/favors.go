package world

import "pure-game-kit/gui/field"

func (ws *WorldScreen) handleFavorsPopup() {
	// cool effect idea - pressing the + should make the whole favor slide to the left
	// (towards the notebook) before updating all favors in the popup
	ws.favors.SetField("title-bgr", field.Text, "Favors to Grant")
	ws.tryExitPopup(ws.favors, ws.settlement, nil)
}
