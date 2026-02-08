package world

import (
	"pure-game-kit/data/file"
	"pure-game-kit/data/path"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/text"
)

var referencedPopups [5]string

func (ws *WorldScreen) handleEventsPopup() {
	var _, textH = ws.events.AreaText("text", ws.camera)
	var enter = keyboard.IsKeyJustPressed(key.Enter) || keyboard.IsKeyJustPressed(key.NumpadEnter)

	ws.events.SetField("text", field.Height, text.New(textH))

	for i := range 5 {
		if ws.events.IsButtonJustClicked(text.New("choice", i+1), ws.camera) {
			if referencedPopups[i] != "" {
				ws.loadEventFile(referencedPopups[i])
				break
			}

			ws.currentPopup = nil
			break
		}
	}

	if enter && ws.events.InputFieldTyping() == "file-name" {
		var fileName = ws.events.Field("file-name", field.Text, ws.camera)
		ws.loadEventFile(fileName)
	}
}
func (ws *WorldScreen) loadEventFile(fileName string) {
	var filePath = path.New("data", "events", fileName+".txt")

	if !file.Exists(filePath) {
		return
	}

	var content = text.Trim(file.LoadText(filePath))
	var lines = text.SplitLines(content)
	var choiceAmount = 0
	var story = ""

	referencedPopups = [5]string{}
	ws.events.InputFieldStopTyping()
	ws.events.SetField("file-name", field.Text, "")
	ws.events.SetField("title-label", field.Text, fileName)

	for i := range 5 {
		ws.events.SetField(text.New("choice", i+1), field.Hidden, "1")
	}

	for _, line := range lines {
		if choiceAmount < 5 && (text.StartsWith(line, "...") || text.StartsWith(line, "…")) {
			choiceAmount++
			var id = text.New("choice", choiceAmount)
			ws.events.SetField(id, field.Text, line)
			ws.events.SetField(id, field.Hidden, "")
			continue
		}

		if text.StartsWith(line, ">>> event ") {
			referencedPopups[choiceAmount-1] = text.Remove(line, ">>> event ")
		}

		var firstSymbol = text.Part(line, 0, 1)
		var quoteOrEmptyLine = firstSymbol == "\"" || firstSymbol == "“" || line == ""
		var digitOrLetter = text.IsAllDigits(firstSymbol) || text.IsAllLetters(firstSymbol)
		if quoteOrEmptyLine || digitOrLetter {
			story += line + "\n"
		}
	}
	ws.events.SetField("text", field.Text, text.Trim(story))
	ws.events.SetField("event", field.Height, text.New(345+(5-choiceAmount)*55))
}
