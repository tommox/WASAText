package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) manageGroupUsersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGroupUsers: invalid Group_id")
		return
	}

	userId, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGroupUsers: invalid User_id")
		return
	}

	requestingUserIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("manageGroupUsers: unauthorized user")
		return
	}

	requestingUserId, err := strconv.Atoi(requestingUserIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGroupUsers: invalid admin ID")
		return
	}

	state := r.URL.Query().Get("state")
	if state != "add" && state != "remove" && state != "promote" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("manageGroupUsers: invalid state parameter")
		return
	}

	if state == "remove" && requestingUserId == userId {
		err = rt.db.RemoveUserFromGroup(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error removing user from group")
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User removed themselves from the group successfully",
		}); err != nil {
			ctx.Logger.WithError(err).Error("manageGroupUsers: error encoding response (self remove)")
		}
		return
	}

	isAdmin, err := rt.db.IsGroupAdmin(groupId, requestingUserId)
	if err != nil || !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("manageGroupUsers: user is not an admin")
		return
	}

	if state == "promote" {
		err = rt.db.PromoteToAdmin(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error promoting user to admin")
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"group_id": groupId,
			"user_id":  userId,
			"role":     "admin",
			"message":  "User promoted to admin successfully",
		}); err != nil {
			ctx.Logger.WithError(err).Error("manageGroupUsers: error encoding response (promote)")
		}
	} else if state == "add" {
		var body struct {
			Role string `json:"role"`
		}
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.Role == "" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("manageGroupUsers: invalid request body for add")
			return
		}

		err = rt.db.AddUserToGroup(groupId, userId, body.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error adding user to group")
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"group_id": groupId,
			"user_id":  userId,
			"role":     body.Role,
			"message":  "User added to group successfully",
		}); err != nil {
			ctx.Logger.WithError(err).Error("manageGroupUsers: error encoding response (add)")
		}
	} else if state == "remove" {
		err = rt.db.RemoveUserFromGroup(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error removing user from group")
			return
		}

		err = rt.db.DeleteAllMessagesFromUserInGroup(groupId, userId)
		if err != nil {
			ctx.Logger.WithError(err).Error("manageGroupUsers: errore nella cancellazione dei messaggi dell'utente")
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User removed from group successfully",
		}); err != nil {
			ctx.Logger.WithError(err).Error("manageGroupUsers: error encoding response (remove)")
		}
	}
}
