package idgenerator

import (
	"context"

	"github.com/fsufitch/bounce-paste/common"
	"github.com/fsufitch/bounce-paste/mq"
)

type ServerContext context.Context

type Server struct {
	Context ServerContext
	Log     common.Logger
	MQ      mq.MQ
}

func (s Server) Run() error {
	s.Log.Warningf("have a good day! %s %v", "foo", s.MQ)

	s.MQ.ManageConnectionAsync(s.Context)

	s.Log.Infof("trying to get MQ session")

	mqSession, err := s.MQ.GetSession(s.Context)
	if err != nil {
		s.Log.Errorf("failed to get MQ session: %v", err)
	} else {
		s.Log.Infof("got MQ connection! %v", &mqSession)
	}

	<-s.Context.Done()

	s.Log.Infof("server shutting down")

	return nil
}
