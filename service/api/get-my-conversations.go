package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getMyConversationsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai `User_id` dal token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("getMyConversations: no valid token")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getMyConversations: invalid user ID")
		return
	}

	// Ottieni le conversazioni
	conversations, err := rt.db.GetUserConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getMyConversations: failed to retrieve conversations")
		return
	}

	// Risposta con successo
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		ctx.Logger.WithError(err).Error("getMyConversations: errore durante l'encoding JSON")
	}
}
