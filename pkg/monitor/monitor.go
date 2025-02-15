package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
)

type Monitor struct {
	db            *database.Queries
	cfg           *config.Config
	monitorTicker *time.Ticker
	monitorCtx    context.Context
	stopMonitor   context.CancelFunc
	stoppedChan   chan struct{}
}

func (m *Monitor) run() {
	for {
		select {
		case <-m.monitorTicker.C:
			slog.Info("Monitor tick")
		case <-m.monitorCtx.Done():
			close(m.stoppedChan)
			return
		}
	}
}

func (m *Monitor) Start() {
	go m.run()
}

func (m *Monitor) Stop() {
	m.stopMonitor()
	<-m.stoppedChan
}

func New(db *database.Queries, cfg *config.Config) *Monitor {
	ctx := context.Background()
	monitorCtx, stopMonitor := context.WithCancel(ctx)

	return &Monitor{
		db:            db,
		cfg:           cfg,
		monitorTicker: time.NewTicker(time.Second * 30),
		monitorCtx:    monitorCtx,
		stopMonitor:   stopMonitor,
		stoppedChan:   make(chan struct{}),
	}
}
