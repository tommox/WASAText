package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) changeGroupNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("changeGroupName: invalid Group_id")
		return
	}

	// Estrai l'utente loggato dal Bearer Token
	adminIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("changeGroupName: unauthorized user")
		return
	}

	adminId, err := strconv.Atoi(adminIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("changeGroupName: invalid admin ID")
		return
	}

	// Controlla che l'utente sia admin o creatore del gruppo
	isAdmin, err := rt.db.IsGroupAdmin(groupId, adminId)
	if err != nil || !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("changeGroupName: user is not an admin")
		return
	}

	// Decodifica il corpo della richiesta
	var body struct {
		GroupName string `json:"group_name"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.GroupName == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("changeGroupName: invalid request body")
		return
	}

	// Aggiorna il nome del gruppo
	err = rt.db.ChangeGroupName(groupId, body.GroupName)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("changeGroupName: group not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("changeGroupName: error updating group name")
		}
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id":   groupId,
		"group_name": body.GroupName,
		"message":    "Group name updated successfully",
	})
}
