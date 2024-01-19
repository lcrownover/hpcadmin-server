package api

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
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
	if u.Name == "" {
		return fmt.Errorf("missing required pirg name: %+v", u)
	}
	if u.OwnerId == 0 {
		return fmt.Errorf("missing required pirg owner_id: %+v", u)
	}
	// owner_id must be present in both admin_ids and user_ids
	if !slices.Contains(u.AdminIds, u.OwnerId) {
		return fmt.Errorf("owner_id %d must be present in admin_ids: %+v", u.OwnerId, u)
	}
	if !slices.Contains(u.UserIds, u.OwnerId) {
		return fmt.Errorf("owner_id %d must be present in user_ids: %+v", u.OwnerId, u)
	}
	// admin_ids must be a subset of user_ids
	for _, adminId := range u.AdminIds {
		if !slices.Contains(u.UserIds, adminId) {
			return fmt.Errorf("admin_id %d must be a present in user_ids: %+v", adminId, u)
		}
	}
	// pirg name must be alphanumeric, lowercase, and start with a letter
	if err := ValidatePirgName(u.Name); err != nil {
		return fmt.Errorf("invalid pirg name: %+v, error: %v", u, err)
	}

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
	slog.Debug("getting all pirgs", "package", "api", "method", "GetAllPirgs")
	var pirgs []*data.Pirg

	pirgs, err := data.GetAllPirgs(h.dbConn)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	resp := newPirgResponseList(pirgs)
	if err := render.RenderList(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreatePirg creates a new Pirg
func (h *PirgHandler) CreatePirg(w http.ResponseWriter, r *http.Request) {
	slog.Debug("creating new pirg", "package", "api", "method", "CreatePirg")
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

	resp := newPirgResponse(newPirg)
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

		pirgIDParam := chi.URLParam(r, "pirgID")
		slog.Debug("loading specific pirg ctx", "id", pirgIDParam, "package", "api", "method", "UserCtx")
		pirgId, err := strconv.Atoi(pirgIDParam)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		pirg, err = data.GetPirgById(h.dbConn, pirgId)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), keys.PirgKey, pirg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPirg returns the Pirg by the ID in the URL
func (h *PirgHandler) GetPirg(w http.ResponseWriter, r *http.Request) {
	slog.Debug("getting pirg", "package", "api", "method", "GetPirg")
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	resp := newPirgResponse(pirg)
	if err := render.Render(w, r, resp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// UpdatePirg updates a Pirg
func (h *PirgHandler) UpdatePirg(w http.ResponseWriter, r *http.Request) {
	slog.Debug("updating pirg", "package", "api", "method", "UpdatePirg")
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	pirgReq := newPirgRequest(pirg)
	if err := render.Bind(r, pirgReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	dataPirgRequest := data.PirgRequest(*pirgReq)
	fmt.Printf("dataPirgRequest: %+v\n", dataPirgRequest)
	updatedPirg, err := data.UpdatePirg(h.dbConn, pirg.Id, &dataPirgRequest)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	resp := newPirgResponse(updatedPirg)
	render.Status(r, http.StatusOK)
	render.Render(w, r, resp)
}

// DeletePirg deletes a Pirg
func (h *PirgHandler) DeletePirg(w http.ResponseWriter, r *http.Request) {
	slog.Debug("deleting pirg", "package", "api", "method", "DeletePirg")
	pirg := r.Context().Value(keys.PirgKey).(*data.Pirg)
	err := data.DeletePirg(h.dbConn, pirg.Id)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.Status(r, http.StatusNoContent)
}

// Utilities
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}
func IsLower(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}
func StartsWithLetter(s string) bool {
	c := s[0]
	if c >= 'a' && c <= 'z' {
		return true
	}
	return false
}
func ValidatePirgName(name string) error {
	errStr := "pirg name must be at least 8 characters, alphanumeric, lowercase, and start with a letter"
	if !IsAlphaNumeric(name) {
		return fmt.Errorf(errStr)
	}
	if !IsLower(name) {
		return fmt.Errorf(errStr)
	}
	if !StartsWithLetter(name) {
		return fmt.Errorf(errStr)
	}
	return nil
}
