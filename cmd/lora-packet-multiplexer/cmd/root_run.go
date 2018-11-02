package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/brocaar/lora-packet-multiplexer/internal/config"
	"github.com/brocaar/lora-packet-multiplexer/internal/multiplexer"
)

func run(cmd *cobra.Command, args []string) error {
	m, err := multiplexer.New(config.C.PacketMultiplexer)
	if err != nil {
		return errors.Wrap(err, "new multiplexer error")
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping lora-packet-multiplexer")
		if err := m.Close(); err != nil {
			log.Fatal(err)
		}
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}
