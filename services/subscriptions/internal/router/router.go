package router

import (
	"net/http"
	"subscriptions/internal/middleware"
)

type Handlers interface {
	CreateSub(w http.ResponseWriter, r *http.Request)
	ReadSub(w http.ResponseWriter, r *http.Request)
	ReadSubs(w http.ResponseWriter, r *http.Request)
	UpdateSub(w http.ResponseWriter, r *http.Request)
	DeleteSub(w http.ResponseWriter, r *http.Request)
	ShowSubscSum(w http.ResponseWriter, r *http.Request)
}
type Router struct {
	r Handlers
}

func NewRouter(a Handlers) *Router {
	return &Router{
		r: a,
	}
}

func (router *Router) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			router.r.ReadSubs(w, r)
		case http.MethodPost:
			router.r.CreateSub(w, r)
		default:
			http.Error(w, "this method are not allowed on this path", http.StatusForbidden)
		}
	})

	mux.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			router.r.DeleteSub(w, r)
		case http.MethodPut:
			router.r.UpdateSub(w, r)
		case http.MethodGet:
			router.r.ReadSub(w, r)
		default:
			http.Error(w, "this method are not allowed on this path", http.StatusForbidden)
		}
	})
	mux.HandleFunc("/subscriptions/sum/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			router.r.ShowSubscSum(w, r)
		default:
			http.Error(w, "this method are not allowed on this path", http.StatusForbidden)
		}
	})
}

func (router *Router) WrapMiddle(mux *http.ServeMux) http.Handler {
	finalmux := middleware.Middleware(mux)
	return finalmux
}
