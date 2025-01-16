package shutdown

import (
	"os"
	"os/signal"
)

type Shutdowner interface {
	Shutdown()
}

func ListenForSignals(signals []os.Signal, shutdowners ...Shutdowner) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	<-c

	for _, shutdowner := range shutdowners {
		shutdowner.Shutdown()
	}

	os.Exit(0)
}
