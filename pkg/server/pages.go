package server

import (
	"net/http"

	"github.com/nint8835/interruption-spotter/pkg/server/ui/pages"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	pages.Index().Render(r.Context(), w)
}
