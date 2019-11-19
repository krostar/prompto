package color

type Config struct {
	Foreground *uint8 `yaml:"fg"`
	Background *uint8 `yaml:"bg"`
}

func (c *Config) SetDefaultColor(def Config) {
	if c.Foreground == nil {
		c.Foreground = def.Foreground
	}

	if c.Background == nil {
		c.Background = def.Background
	}
}

func (cfg Config) ToStyle() Style {
	var (
		fg Color
		bg Color
	)

	if cfg.Foreground != nil {
		fg = NewFGColor(*cfg.Foreground)
	}

	if cfg.Background != nil {
		bg = NewBGColor(*cfg.Background)
	}

	return NewStyle(fg, bg)
}
