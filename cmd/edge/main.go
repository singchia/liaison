package main

import (
	"context"

	"github.com/jumboframes/armorigo/log"
	"github.com/jumboframes/armorigo/sigaction"
	"github.com/singchia/liaison/pkg/edge"
)

func main() {
	edge, err := edge.NewEdge()
	if err != nil {
		log.Errorf("new edge err: %s", err)
		return
	}

	sig := sigaction.NewSignal()
	sig.Wait(context.TODO())

	edge.Close()
}
