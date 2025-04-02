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
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("deleteGroup: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteGroup: invalid user ID")
		return
	}

	isMember, err := rt.db.IsGroupMember(groupId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("addToGroup: error checking group membership")
		return
	}

	if !isMember {
		// Se l'utente non è membro del gruppo, restituiamo un errore 403 Forbidden
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user is not a member of group")).Error("permission denied")
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Group deleted successfully",
	}); err != nil {
		ctx.Logger.WithError(err).Error("deleteGroup: errore durante l'encoding della risposta JSON")
	}
}
