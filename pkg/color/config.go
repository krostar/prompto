package color

// Config stores color configuration.
type Config struct {
	Foreground *string `yaml:"fg"`
	Background *string `yaml:"bg"`
}

// SetDefaultColor sets default to unsetted values.
func (c *Config) SetDefaultColor(def Config) {
	if c.Foreground == nil {
		c.Foreground = def.Foreground
	}

	if c.Background == nil {
		c.Background = def.Background
	}
}

// ToStyle creates a Style from Config.
func (c Config) ToStyle() Style {
	var (
		fg Color
		bg Color
	)

	if c.Foreground != nil {
		fg = NewHexFGColor(*c.Foreground)
	}

	if c.Background != nil {
		bg = NewHexBGColor(*c.Background)
	}

	return NewStyle(fg, bg)
}
