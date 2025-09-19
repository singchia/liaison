package main

import (
	"context"

	"github.com/jumboframes/armorigo/sigaction"
	"github.com/singchia/liaison/pkg/lerrors"
	"github.com/singchia/liaison/pkg/liaison"
	"github.com/sirupsen/logrus"
)

func main() {
	liaison, err := liaison.NewLiaison()
	if err != nil {
		if err != lerrors.ErrInvalidUsage {
			logrus.Errorf("new liaison err: %s", err)
		}
		return
	}

	go liaison.Serve()

	sig := sigaction.NewSignal()
	sig.Wait(context.TODO())

	liaison.Close()
}
