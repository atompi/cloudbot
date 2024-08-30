package utils

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func GracefulExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	zap.S().Warnf("a %v signal is received, exiting...", s)
	os.Exit(0)
}
