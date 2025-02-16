package server

import (
	"net/http"

	"github.com/nint8835/interruption-spotter/pkg/server/ui/pages"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	regions, err := s.db.GetRegions(r.Context())
	if err != nil {
		s.logger.Error("Failed to get regions", "err", err)
		http.Error(w, "failed to get regions", http.StatusInternalServerError)
		return
	}

	instanceTypes, err := s.db.GetInstanceTypes(r.Context())
	if err != nil {
		s.logger.Error("Failed to get instance types", "err", err)
		http.Error(w, "failed to get instance types", http.StatusInternalServerError)
		return
	}

	operatingSystems, err := s.db.GetOperatingSystems(r.Context())
	if err != nil {
		s.logger.Error("Failed to get operating systems", "err", err)
		http.Error(w, "failed to get operating systems", http.StatusInternalServerError)
		return
	}

	err = pages.Index(pages.IndexProps{
		Regions:          regions,
		InstanceTypes:    instanceTypes,
		OperatingSystems: operatingSystems,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.Error("Failed to render index page", "err", err)
		return
	}
}
