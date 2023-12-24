package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/data"
	"github.com/lcrownover/hpcadmin-server/internal/keys"
)

type PirgResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerId   int       `json:"owner_id"`
	AdminIds  []int     `json:"admin_ids"`
	UserIds   []int     `json:"user_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *PirgResponse) Bind(r *http.Request) error {
	return nil
}

func (u *PirgResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newPirgResponse(u *data.Pirg) *PirgResponse {
	return &PirgResponse{
		Id:        u.Id,
		Name:      u.Name,
		OwnerId:   u.OwnerId,
		AdminIds:  u.AdminIds,
		UserIds:   u.UserIds,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// newPirgResponseList converts a list of PirgResponse objects into a list of render.Renderer objects
func newPirgResponseList(Pirgs []*data.Pirg) []render.Renderer {
	list := []render.Renderer{}
	for _, pirg := range Pirgs {
		list = append(list, newPirgResponse(pirg))
	}
	return list
}

type PirgRequest struct {
	Name     string `json:"name"`
	OwnerId  int    `json:"owner_id"`
	AdminIds []int  `json:"admin_ids"`
	UserIds  []int  `json:"user_ids"`
}

func (u *PirgRequest) Bind(r *http.Request) error {
	if u.Name == "" || u.OwnerId == 0 {
		return fmt.Errorf("missing required Pirg fields: %+v", u)
	}
	// add in more checks like alphanumeric, length, etc.
	return nil
}

func newPirgRequest(u *data.Pirg) *PirgRequest {
	return &PirgRequest{
		Name:     u.Name,
		OwnerId:  u.OwnerId,
		AdminIds: u.AdminIds,
		UserIds:  u.UserIds,
	}
}

type PirgStub struct {
	Id       int
	Pirgname string
}

type PirgHandler struct {
	dbConn *sql.DB
}

func PirgsRouter(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	h := newPirgHandler(ctx)
	r.Get("/", h.GetAllPirgs)
	r.Post("/", h.CreatePirg)
	r.Route("/{pirgID}", func(r chi.Router) {
		r.Use(h.PirgCtx)
		r.Get("/", h.GetPirg)
		r.Put("/", h.UpdatePirg)
		r.Delete("/", h.DeletePirg)
		// r.Mount("/admins", PirgAdminsRouter(ctx))
	})
	return r
}

func newPirgHandler(ctx context.Context) *PirgHandler {
	dbConn := ctx.Value(keys.DBConnKey).(*sql.DB)
	return &PirgHandler{dbConn: dbConn}
}

// GetAllPirgs returns all existing Pirgs
func (h *PirgHandler) GetAllPirgs(w http.ResponseWriter, r *http.Request) {
	resp := &ApiResponse{}
	var pirgs []*data.Pirg

	pirgs, err := data.GetAllPirgs(h.dbConn)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	resp.Results = newPirgResponseList(pirgs)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// GetPirgById returns a single Pirg by id, but is not currently used
func (h *PirgHandler) GetPirgById(w http.ResponseWriter, r *http.Request) {
	resp := &ApiResponse{}
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)

	resp.Results = newPirgResponse(pirg)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreatePirg creates a new Pirg
func (h *PirgHandler) CreatePirg(w http.ResponseWriter, r *http.Request) {
	resp := &ApiResponse{}
	pirg := &PirgRequest{}
	if err := render.Bind(r, pirg); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	dataPirg := data.PirgRequest(*pirg)

	newPirg, err := data.CreatePirg(h.dbConn, &dataPirg)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp.Results = newPirgResponse(newPirg)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, resp)
}

// PirgCtx middleware is used to load a Pirg object from /Pirgs/{Pirgname} requests
// and then attach it to the request context. In case of failure the request is aborted
// and a 404 error response is sent to the client.
func (h *PirgHandler) PirgCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pirg *data.Pirg
		var err error

		pirgIDParam := chi.URLParam(r, "PirgID")
		PirgId, err := strconv.Atoi(pirgIDParam)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		pirg, err = data.GetPirgById(h.dbConn, PirgId)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keys.PirgKey, pirg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPirg returns the Pirg in the request context
func (h *PirgHandler) GetPirg(w http.ResponseWriter, r *http.Request) {
	resp := &ApiResponse{}
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	resp.Results = newPirgResponse(pirg)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// UpdatePirg updates a Pirg
func (h *PirgHandler) UpdatePirg(w http.ResponseWriter, r *http.Request) {
	resp := &ApiResponse{}
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	pirgReq := newPirgRequest(pirg)
	if err := render.Bind(r, pirgReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	dataPirgRequest := data.PirgRequest(*pirgReq)
	updatedPirg, err := data.UpdatePirg(h.dbConn, pirg.Id, &dataPirgRequest)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	resp.Results = newPirgResponse(updatedPirg)
	render.Status(r, http.StatusOK)
	render.Render(w, r, resp)
}

// DeletePirg deletes a Pirg
func (h *PirgHandler) DeletePirg(w http.ResponseWriter, r *http.Request) {
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	err := data.DeletePirg(h.dbConn, pirg.Id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.Status(r, http.StatusNoContent)
}
