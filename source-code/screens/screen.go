package screens

type Screen interface {
	OnLoad()
	OnUpdate()
}

const Menu, World, Battle = 0, 1, 2

var current = 1
var all []Screen

func New(screens ...Screen) {
	all = screens
}

func Current() Screen {
	return all[current]
}
func UpdateCurrent() {
	all[current].OnUpdate()
}
func LoadAll() {
	for _, scr := range all {
		if scr != nil {
			scr.OnLoad()
		}
	}
}
