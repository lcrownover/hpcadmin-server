package api

import "net/http"

type ApiResponse struct {
	Results any `json:"results"`
}

func (u *ApiResponse) Bind(r *http.Request) error {
	return nil
}

func (u *ApiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
