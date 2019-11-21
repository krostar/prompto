package segment

import (
	"errors"
	usr "os/user"

	"github.com/krostar/prompto/pkg/color"

	"github.com/krostar/prompto/internal/prompto/domain"
)

type user struct {
	cfg userConfig
}

type userConfig struct {
	Color color.Config `yaml:"color"`
}

func segmentUser(rcfg interface{}) (domain.SegmentsProvider, error) {
	cfg, isArgConfig := rcfg.(userConfig)
	if !isArgConfig {
		return nil, errors.New("segmentuser expected 1 arg of type userConfig")
	}

	return &user{cfg: cfg}, nil
}

func (u *user) SegmentName() string {
	return "user"
}

func (u *user) ProvideSegments() (domain.Segments, error) {
	cu, err := usr.Current()
	if err != nil {
		return nil, err
	}
	return domain.Segments{
		domain.
			NewSegment(cu.Username).
			WithSpaceAround().
			SetStyle(u.cfg.Color.ToStyle()),
	}, nil
}
