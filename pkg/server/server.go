package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/feeds"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
)

type Server struct {
	cfg *config.Config
	db  *database.Queries
	mux *http.ServeMux
}

func (s *Server) Run() error {
	if err := http.ListenAndServe(s.cfg.BindAddr, s.mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if !(query.Has("regions") && query.Has("instance_types") && query.Has("operating_systems")) {
		http.Error(w, "missing required query parameters", http.StatusBadRequest)
		return
	}

	regions := strings.Split(query.Get("regions"), ",")
	instanceTypes := strings.Split(query.Get("instance_types"), ",")
	operatingSystems := strings.Split(query.Get("operating_systems"), ",")

	changes, err := s.db.GetInterruptionChanges(r.Context(), database.GetInterruptionChangesParams{
		Regions:          regions,
		InstanceTypes:    instanceTypes,
		OperatingSystems: operatingSystems,
	})
	if err != nil {
		slog.Error("Failed to get interruption changes", "err", err)
		http.Error(w, "failed to get interruption changes", http.StatusInternalServerError)
		return
	}

	feed := &feeds.Feed{
		Title: "Interruption Spotter - Spot Advisor Changes",
		Link: &feeds.Link{
			Href: "https://interruption-spotter.bootleg.technology",
		},
		Description: fmt.Sprintf("Spot Advisor changes for regions %s, instance types %s, and operating systems %s", regions, instanceTypes, operatingSystems),
		Author: &feeds.Author{
			Name:  "Interruption Spotter",
			Email: "interruption-spotter@rileyflynn.me",
		},
		Created: time.Now(),
		Updated: changes[0].ObservedTime,
	}

	for _, change := range changes {
		item := &feeds.Item{
			Id:      fmt.Sprintf("%d", change.ID),
			Created: change.ObservedTime,
			Title:   fmt.Sprintf("Spot Advisor Change for %s on %s in %s", change.InstanceType, change.OperatingSystem, change.Region),
		}

		prevLevel, hasPrevLevel := change.LastInterruptionLevel.(int64)
		if !hasPrevLevel {
			item.Description = fmt.Sprintf(
				"Instance type %s in region %s running %s now has interruption level %d (%s).",
				change.InstanceType,
				change.Region,
				change.OperatingSystem,
				change.InterruptionLevel,
				change.InterruptionLevelLabel,
			)
		} else {
			prevLevelLabel, _ := change.LastInterruptionLevelLabel.(string)
			item.Description = fmt.Sprintf(
				"Instance type %s in region %s running %s changed from interruption level %d (%s) to %d (%s).",
				change.InstanceType,
				change.Region,
				change.OperatingSystem,
				prevLevel,
				prevLevelLabel,
				change.InterruptionLevel,
				change.InterruptionLevelLabel,
			)
		}

		feed.Add(item)
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	err = feed.WriteRss(w)
	if err != nil {
		slog.Error("Failed to write RSS feed", "err", err)
		http.Error(w, "failed to write RSS feed", http.StatusInternalServerError)
	}
}

func New(cfg *config.Config, db *database.Queries) (*Server, error) {
	instance := &Server{
		cfg: cfg,
		db:  db,
		mux: http.NewServeMux(),
	}

	instance.mux.HandleFunc("/feed", instance.handleFeed)

	return instance, nil
}
