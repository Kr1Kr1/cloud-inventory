package utils

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"
)

var NopShutdown = NewGroupShutdown()

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type GroupShutdown struct {
	s []Shutdowner
}

func NewGroupShutdown(s ...Shutdowner) *GroupShutdown {
	return &GroupShutdown{s}
}

func (g *GroupShutdown) Shutdown(ctx context.Context) error {
	var errm error
	for _, s := range g.s {
		if err := s.Shutdown(ctx); err != nil {
			errm = multierror.Append(errm, err)
		}
	}
	return errm
}
