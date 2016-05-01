package acme

import (
	"github.com/tiagofalcao/gogae_common/engine"
	"net/http"
)

func init() {
	http.HandleFunc("/.well-known/acme-challenge", engine.HandleAdminReq(index))
	http.HandleFunc("/.well-known/acme-challenge-renew", engine.HandleAdminReq(renew))
	http.HandleFunc("/.well-known/acme-challenge/", engine.Handle(response))
}
