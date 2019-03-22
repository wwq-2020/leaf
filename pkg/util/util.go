package util

import (
	"os"
	"os/signal"
	"syscall"
)

func OnExit(cleanups ...func()) {
	go onExit(cleanups...)
}

func onExit(cleanups ...func()) {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	for _, cleanup := range cleanups {
		cleanup()
	}
}
