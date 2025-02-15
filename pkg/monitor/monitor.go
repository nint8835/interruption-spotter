package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
	"github.com/nint8835/interruption-spotter/pkg/spotdata"
)

type Monitor struct {
	db            *database.Queries
	cfg           *config.Config
	fetcher       *spotdata.Fetcher
	monitorTicker *time.Ticker
	monitorCtx    context.Context
	stopMonitor   context.CancelFunc
	stoppedChan   chan struct{}
}

func (m *Monitor) updateInterruptionLevels() error {
	slog.Debug("Checking if we should fetch spot data")

	shouldFetch, err := m.fetcher.ShouldFetch(m.monitorCtx)
	if err != nil {
		return fmt.Errorf("failed to check if we should fetch spot data: %w", err)
	}

	if !shouldFetch {
		slog.Debug("Spot data file unchanged, skipping fetch")
		return nil
	}

	slog.Debug("Fetching spot data")
	spotData, err := m.fetcher.Fetch(m.monitorCtx)
	if err != nil {
		return fmt.Errorf("failed to fetch spot data: %w", err)
	}

	currentVals, err := m.db.GetCurrentInterruptionLevels(m.monitorCtx)
	if err != nil {
		return fmt.Errorf("failed to get current interruption levels: %w", err)
	}

	type currentValKey struct {
		region          string
		operatingSystem string
		instanceType    string
	}
	currentValMap := make(map[currentValKey]database.GetCurrentInterruptionLevelsRow)
	for _, val := range currentVals {
		currentValMap[currentValKey{
			region:          val.Region,
			operatingSystem: val.OperatingSystem,
			instanceType:    val.InstanceType,
		}] = val
	}

	for regionName, regionStats := range spotData.SpotAdvisor {
		for osName, osStats := range regionStats {
			for instanceType, instanceStats := range osStats {
				currentVal, hasCurrentVal := currentValMap[currentValKey{
					region:          regionName,
					operatingSystem: osName,
					instanceType:    instanceType,
				}]
				interruptionLevelLabel := spotData.Ranges[instanceStats.InterruptionLevel].Label

				if hasCurrentVal &&
					currentVal.InterruptionLevel == instanceStats.InterruptionLevel &&
					currentVal.InterruptionLevelLabel == interruptionLevelLabel {
					slog.Debug(
						"Interruption level unchanged",
						"region", regionName,
						"os", osName,
						"instance_type", instanceType,
						"interruption_level", instanceStats.InterruptionLevel,
						"interruption_level_label", interruptionLevelLabel,
					)
					continue
				}

				slog.Info(
					"Interruption level changed",
					"region", regionName,
					"os", osName,
					"instance_type", instanceType,
					"interruption_level", instanceStats.InterruptionLevel,
					"interruption_level_label", interruptionLevelLabel,
				)
				err = m.db.InsertStat(m.monitorCtx, database.InsertStatParams{
					Region:                 regionName,
					OperatingSystem:        osName,
					InstanceType:           instanceType,
					InterruptionLevel:      instanceStats.InterruptionLevel,
					InterruptionLevelLabel: interruptionLevelLabel,
				})
				if err != nil {
					return fmt.Errorf("failed to insert stat: %w", err)
				}
			}
		}
	}

	slog.Debug("Finished updating interruption levels")

	return nil
}

func (m *Monitor) run() {
	for {
		select {
		case <-m.monitorTicker.C:
			err := m.updateInterruptionLevels()
			if err != nil {
				slog.Error("Failed to update interruption levels", "err", err)
			}
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
		fetcher:       &spotdata.Fetcher{},
		monitorTicker: time.NewTicker(cfg.PollInterval),
		monitorCtx:    monitorCtx,
		stopMonitor:   stopMonitor,
		stoppedChan:   make(chan struct{}),
	}
}
