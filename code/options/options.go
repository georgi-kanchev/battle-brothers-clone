package options

type Options struct {
	Graphics struct {
		WindowState int  `yaml:"window-state"`
		Monitor     int  `yaml:"monitor"`
		VSync       bool `yaml:"vsync"`
		LimitFPS    int  `yaml:"limit-fps"`
	} `yaml:"graphics"`
	ScaleUI struct {
		Master float32 `yaml:"master"`
		Menu   struct {
			Options float32 `yaml:"options"`
		} `yaml:"menu"`
		World struct {
			HUD        float32 `yaml:"hud"`
			Inventory  float32 `yaml:"inventory"`
			Settlement float32 `yaml:"settlement"`
		} `yaml:"world"`
		Battle struct {
			HUD  float32 `yaml:"hud"`
			Loot float32 `yaml:"loot"`
		} `yaml:"battle"`
	} `yaml:"user-interface-scale"`
	AudioVolume struct {
		Master float32 `yaml:"master"`
		Music  float32 `yaml:"music"`
		Sound  float32 `yaml:"sound"`
	} `yaml:"audio-volume"`
	Controls struct {
	} `yaml:"controls"`
	Gameplay struct {
	} `yaml:"gameplay"`
}
