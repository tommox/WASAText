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

func (rt *_router) deleteGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteGroup: invalid Group_id")
		return
	}

	// Estrai l'utente loggato dal Bearer Token
	adminIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("deleteGroup: unauthorized user")
		return
	}

	adminId, err := strconv.Atoi(adminIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteGroup: invalid admin ID")
		return
	}

	// Verifica che l'utente loggato sia admin o creator del gruppo
	isAdmin, err := rt.db.IsGroupAdmin(groupId, adminId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("deleteGroup: error checking admin permissions")
		return
	}

	if !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user not authorized to delete group")).Error("deleteGroup: permission denied")
		return
	}

	// Elimina il gruppo
	err = rt.db.DeleteGroup(groupId)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deleteGroup: group not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("deleteGroup: error deleting group")
		}
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Group deleted successfully",
	})
}
