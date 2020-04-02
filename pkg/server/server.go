package server

import (
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"github.com/danielchen115/mock-oauth2/pkg/server/handler"
	"net/http"
)

func New(svc oauth.Service, addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/authorize", handler.Authorize(svc))
	mux.Handle("/current-user", handler.CurrentUser(svc))
	return &http.Server{Handler: mux, Addr: addr}
}