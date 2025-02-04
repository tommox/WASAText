package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	ctx.Logger.Infof("Received request to get online users from: %s", r.RemoteAddr)

	// Recupera la lista degli utenti online dal database
	users, err := rt.db.GetUsers()
	if err != nil {
		ctx.Logger.WithError(err).Error("getUsers: error fetching online users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Rispondi con la lista degli utenti online
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
