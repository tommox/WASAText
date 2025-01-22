package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) createGroupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai l'utente corrente dal Bearer Token
	creatorIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("createGroup: unauthorized user")
		return
	}

	creatorId, err := strconv.Atoi(creatorIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("createGroup: invalid creator ID")
		return
	}

	// Decodifica il corpo della richiesta
	var body struct {
		GroupName string `json:"group_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.GroupName == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("createGroup: invalid request body")
		return
	}

	// Crea il gruppo nel database
	groupId, err := rt.db.CreateGroup(body.GroupName, creatorId, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("createGroup: error creating group")
		return
	}

	// Rispondi con successo
	response := map[string]interface{}{
		"group_id":   groupId,
		"group_name": body.GroupName,
		"creator_id": creatorId,
		"created_at": time.Now(),
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
