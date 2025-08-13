package serviceprogram

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kardianos/service"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/grpcserver"
)

type Program struct {
	logger   service.Logger
	server   *grpcserver.Server
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	startErr error
	startCh  chan struct{}
}

func New() *Program {
	return &Program{
		startCh: make(chan struct{}),
	}
}

func (p *Program) SetLogger(logger service.Logger) {
	p.logger = logger
}

func (p *Program) Start(s service.Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.server = grpcserver.NewServer()

	p.wg.Add(1)
	go p.run()

	// Wait for server to start or fail
	select {
	case <-p.startCh:
		if p.startErr != nil {
			p.cancel()
			p.wg.Wait()
			return p.startErr
		}
	case <-time.After(5 * time.Second):
		p.cancel()
		p.wg.Wait()
		return fmt.Errorf("timeout waiting for server to start")
	}

	return nil
}

func (p *Program) Stop(s service.Service) error {
	if p.cancel != nil {
		p.cancel()
	}

	if p.server != nil {
		if err := p.server.Stop(); err != nil {
			if p.logger != nil {
				p.logger.Errorf("Error stopping server: %v", err)
			}
		}
	}

	p.wg.Wait()
	return nil
}

func (p *Program) run() {
	defer p.wg.Done()
	defer close(p.startCh)

	grpcAddr := config.DefaultBindAddress + ":" + config.GRPCPort
	grpcWebAddr := config.DefaultBindAddress + ":" + config.GRPCWebPort

	// Start the server
	if err := p.server.Start(grpcAddr, grpcWebAddr); err != nil {
		p.startErr = err
		return
	}

	// Signal successful start
	p.startCh <- struct{}{}

	// Wait a moment for server to be ready
	time.Sleep(2 * time.Second)

	go p.attemptAutoLogin()

	<-p.ctx.Done()
}

func (p *Program) attemptAutoLogin() {
	mnemonic, err := core.GetStoredMnemonic()
	if err != nil || mnemonic == "" {
		if p.logger != nil {
			p.logger.Info("No stored mnemonic found, skipping auto-login")
		}
		return
	}

	if p.logger != nil {
		p.logger.Info("Found stored mnemonic, attempting auto-login...")
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := core.LoginAccount(mnemonic, "", ""); err != nil {
			if i < maxRetries-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			if p.logger != nil {
				p.logger.Errorf("Failed to auto-login after %d attempts: %v", maxRetries, err)
			}
		} else {
			if p.logger != nil {
				p.logger.Info("Successfully logged in using stored mnemonic")
			}
			break
		}
	}
}
