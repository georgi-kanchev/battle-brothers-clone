package world

import "pure-game-kit/gui/field"

func (ws *WorldScreen) handleRecruitPopup() {
	ws.recruit.SetField("title-bgr", field.Text, "Adventurers to Recruit")
	ws.tryExitPopup(ws.recruit, ws.settlement, nil)
}
