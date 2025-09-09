package main

import (
	"context"

	"github.com/jumboframes/armorigo/sigaction"
	"github.com/singchia/liaison/pkg/edge"
	"github.com/sirupsen/logrus"
)

func main() {
	edge, err := edge.NewEdge()
	if err != nil {
		logrus.Errorf("new edge err: %s", err)
		return
	}

	sig := sigaction.NewSignal()
	sig.Wait(context.TODO())

	edge.Close()
}
