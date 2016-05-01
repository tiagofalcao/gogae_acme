package acme

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"

	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/tiagofalcao/gogae_common/engine"
)

type ACMEChallenge struct {
	Key       *datastore.Key `datastore:"-"`
	Challenge string
	Response  string
}

var engineKind = "ACMEChallenge"
var engineKindDef = "@"

func index(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	t, err := template.ParseFiles("templates/base.html", "github.com/tiagofalcao/gogae_acme/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k := datastore.NewKey(
		c,             // appengine.Context
		engineKind,    // Kind
		engineKindDef, // String ID; empty means no string ID
		0,             // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		nil,           // Parent Key; nil means no parent
	)

	var challenge ACMEChallenge
	err = engine.Get(c, k, &challenge)

	if err != nil && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = t.Execute(w, challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renew(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	decoder := json.NewDecoder(r.Body)
	var challenge ACMEChallenge
	err := decoder.Decode(&challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k := datastore.NewKey(
		c,             // appengine.Context
		engineKind,    // Kind
		engineKindDef, // String ID; empty means no string ID
		0,             // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		nil,           // Parent Key; nil means no parent
	)

	err = engine.Put(c, k, &challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var validPath = regexp.MustCompile("^/.well-known/acme-challenge/(.+)$")

func response(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	k := datastore.NewKey(
		c,             // appengine.Context
		engineKind,    // Kind
		engineKindDef, // String ID; empty means no string ID
		0,             // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		nil,           // Parent Key; nil means no parent
	)

	var challenge ACMEChallenge
	err := engine.Get(c, k, &challenge)
	if err != nil || challenge.Challenge != m[1] {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprint(w, challenge.Response)
}
