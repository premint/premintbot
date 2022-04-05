package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	magicClient "github.com/magiclabs/magic-admin-go/client"
	"github.com/magiclabs/magic-admin-go/token"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const authBearer = "Bearer"

// ProvideRouter provides a gorilla mux router
func ProvideRouter(lc fx.Lifecycle, logger *zap.SugaredLogger, magic *magicClient.API) *mux.Router {
	var router = mux.NewRouter()

	router.Use(jsonMiddleware)

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				addr := ":8080"
				logger.Info("Listening on ", addr)

				go http.ListenAndServe(addr, router)

				return nil
			},
		},
	)

	return router
}

// jsonMiddleware makes sure that every response is JSON
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

var Options = ProvideRouter

// magicMiddleware validates the incoming user's auth token
func magicMiddleware(magic *magicClient.API) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate token with Magic
			if !strings.HasPrefix(r.Header.Get("Authorization"), authBearer) {
				fmt.Fprintf(w, "Bearer token is required")
				return
			}

			did := r.Header.Get("Authorization")[len(authBearer)+1:]
			if did == "" {
				fmt.Fprintf(w, "DID token is required")
				return
			}

			tk, err := token.NewToken(did)
			if err != nil {
				fmt.Fprintf(w, "Malformed DID token error: %s", err.Error())
				return
			}

			if err := tk.Validate(); err != nil {
				fmt.Fprintf(w, "DID token failed validation: %s", err.Error())
				return
			}

			userInfo, err := magic.User.GetMetadataByIssuer(tk.GetIssuer())
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err.Error())
				return
			}

			fmt.Println(userInfo.Email)

			next.ServeHTTP(w, r)
		})
	}
}
