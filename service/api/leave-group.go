package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) leaveGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("leaveGroup: invalid Group_id")
		return
	}

	// Estrai l'utente loggato dal Bearer Token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("leaveGroup: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("leaveGroup: invalid user ID")
		return
	}

	// Verifica che l'utente sia membro del gruppo
	isMember, err := rt.db.IsGroupMember(groupId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("leaveGroup: error checking group membership")
		return
	}

	if !isMember {
		// Se l'utente non è membro del gruppo, restituiamo un errore 400
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(errors.New("user is not member of group")).Error("permission denied")
		return
	}

	// Rimuovi l'utente dal gruppo
	err = rt.db.RemoveUserFromGroup(groupId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("leaveGroup: error removing user from group")
		return
	}

	// Opzionale: Se si desidera, cancelliamo i messaggi dell'utente dal gruppo
	err = rt.db.DeleteAllMessagesFromUserInGroup(groupId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("leaveGroup: error deleting messages from user")
	}

	// Risposta di successo
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User left the group successfully",
	}); err != nil {
		ctx.Logger.WithError(err).Error("leaveGroup: error encoding response")
	}
}
