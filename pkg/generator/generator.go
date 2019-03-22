package generator

import (
	"errors"

	"github.com/wwq1988/leaf/pkg/conf"
)

var (
	ErrUnknownMode = errors.New("unknown mode")
)

type Generator interface {
	Generate(string) ([]uint64, error)
	Stop()
}

func New() (Generator, error) {
	mode := conf.GetMode()
	switch mode {
	case conf.ModeSegment:
		return NewSegment()
	case conf.ModeSnowflake:
		return NewSnowflake()
	default:
		return nil, ErrUnknownMode
	}
}
