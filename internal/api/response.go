package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type ApiResponse struct {
	Role    string            `json:"role"`
	Results []render.Renderer `json:"results"`
}

func (u *ApiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
