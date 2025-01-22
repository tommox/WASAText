package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) manageGroupUsersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGroupUsers: invalid Group_id")
		return
	}

	// Estrai lo User_id dal percorso
	userId, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGroupUsers: invalid User_id")
		return
	}

	// Estrai l'admin ID dal token Bearer
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

	// Estrai lo stato dall'URL (add, remove, promote)
	state := r.URL.Query().Get("state")
	if state != "add" && state != "remove" && state != "promote" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("manageGroupUsers: invalid state parameter")
		return
	}

	// Autorizzazione per l'operazione
	if state == "remove" && requestingUserId == userId {
		// Se l'utente sta rimuovendo sé stesso, bypassa il controllo admin
		err = rt.db.RemoveUserFromGroup(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error removing user from group")
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User removed themselves from the group successfully",
		})
		return
	}

	// Controlla i permessi dell'admin per tutte le altre operazioni
	isAdmin, err := rt.db.IsGroupAdmin(groupId, requestingUserId)
	if err != nil || !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("manageGroupUsers: user is not an admin")
		return
	}

	// Gestisci le diverse operazioni
	if state == "promote" {
		// Promuovi l'utente al ruolo di admin
		err = rt.db.PromoteToAdmin(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error promoting user to admin")
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"group_id": groupId,
			"user_id":  userId,
			"role":     "admin",
			"message":  "User promoted to admin successfully",
		})
	} else if state == "add" {
		// Decodifica il corpo della richiesta per ottenere il ruolo
		var body struct {
			Role string `json:"role"`
		}
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.Role == "" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("manageGroupUsers: invalid request body for add")
			return
		}

		// Aggiungi l'utente al gruppo
		err = rt.db.AddUserToGroup(groupId, userId, body.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error adding user to group")
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"group_id": groupId,
			"user_id":  userId,
			"role":     body.Role,
			"message":  "User added to group successfully",
		})
	} else if state == "remove" {
		// Rimuovi l'utente dal gruppo
		err = rt.db.RemoveUserFromGroup(groupId, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageGroupUsers: error removing user from group")
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "User removed from group successfully",
		})
	}
}
