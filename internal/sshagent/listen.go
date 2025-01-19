package sshagent

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	sshagent "golang.org/x/crypto/ssh/agent"
)

func Listen(
	ctx context.Context,
	p string,
	a sshagent.ExtendedAgent,
	logger logrus.FieldLogger,
) error {
	_ = os.Remove(p) // remove previous agent
	err := os.MkdirAll(filepath.Dir(p), 0o777)
	if err != nil {
		return errors.Wrap(err, "failed to create directory for socket")
	}

	logger.
		WithFields(logrus.Fields{
			"socket-path": p,
		}).
		Debugln("starting socket")

	l, err := net.Listen("unix", p)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on unix socket %s", p)
	}

	listen := make(chan tuple)
	continueChan := make(chan struct{})

	go acceptToChan(l, listen, continueChan)

	for {
		select {
		case <-ctx.Done():
			_ = l.Close()
			close(continueChan)

			return nil
		case t := <-listen:
			c, err2 := t.Conn, t.error
			if err2 != nil {
				type temporary interface {
					Temporary() bool
				}
				if tempErr, ok := err.(temporary); ok && tempErr.Temporary() {
					logger.Println("Temporary Accept error, sleeping 1s:", tempErr)
					time.Sleep(1 * time.Second)
					continue
				}

				// termination block
				close(continueChan)

				if errors.Is(err2, io.EOF) {
					_ = l.Close()
					return nil
				}

				return errors.Wrap(err2, "failed to accept connection")
			}

			go func() {
				serveErr := sshagent.ServeAgent(a, c)
				if serveErr != nil && !errors.Is(serveErr, io.EOF) {
					logger.Println("Failed to serve agent:", serveErr)
				}
			}()

			continueChan <- struct{}{}
		}
	}
}

type tuple struct {
	net.Conn
	error
}

func acceptToChan(l net.Listener, out chan<- tuple, continueChan <-chan struct{}) {
	for {
		conn, err := l.Accept()
		out <- tuple{
			Conn:  conn,
			error: err,
		}

		_, ok := <-continueChan
		if !ok {
			return
		}
	}
}
