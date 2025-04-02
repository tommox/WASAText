package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) addToGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addToGroup: invalid Group_id")
		return
	}

	userId, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addToGroup: invalid User_id")
		return
	}

	requestingUserIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("addToGroup: unauthorized user")
		return
	}

	requestingUserId, err := strconv.Atoi(requestingUserIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addToGroup: invalid requester ID")
		return
	}

	// Verifica che l'utente che sta facendo la richiesta sia membro del gruppo
	isMember, err := rt.db.IsGroupMember(groupId, requestingUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("addToGroup: error checking group membership")
		return
	}

	if !isMember {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user is not a member of group")).Error("permission denied")
		return
	}

	// Aggiungi l'utente target al gruppo
	err = rt.db.AddUserToGroup(groupId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("addToGroup: error adding user to group")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id": groupId,
		"user_id":  userId,
		"message":  "User added to group successfully",
	}); err != nil {
		ctx.Logger.WithError(err).Error("addToGroup: error encoding response")
	}
}
