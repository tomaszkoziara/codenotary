package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tomaszkoziara/codenotarybe/accounting"
)

func CreateRouter(accountingService *accounting.Accounting) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.Recoverer,
		middleware.Logger,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		},
	)

	router.Post("/api/v0/accountinginfo", CreateStoreAccountingInfo(accountingService))
	router.Get("/api/v0/accountinginfo", CreateGetListAccountingInfo(accountingService))

	return router
}
