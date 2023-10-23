package idgenerator

import (
	"github.com/fsufitch/bounce-paste/common"
)

type Server struct {
	Log common.Logger
}

func (s Server) Run() error {
	s.Log.Warningf("have a good day! %s", "foo")
	return nil
}
