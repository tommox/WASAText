package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai `Conversation_id` dal percorso
	conversationId, err := strconv.Atoi(ps.ByName("Conversation_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getConversation: invalid Conversation_id")
		return
	}

	// Estrai `User_id` dal token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("getConversation: no valid token")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getConversation: invalid user ID")
		return
	}

	// Verifica i permessi dell'utente
	hasAccess, isGroup, err := rt.db.CheckConversationAccess(userId, conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getConversation: error checking access")
		return
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user has no access")).Error("getConversation: user has no access to this conversation")
		return
	}

	// Recupera i messaggi della conversazione
	var messages interface{}
	if isGroup {
		messages, err = rt.db.GetGroupConversationMessages(conversationId)
	} else {
		messages, err = rt.db.GetConversationMessages(conversationId)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getConversation: failed to retrieve messages")
		return
	}

	// Rispondi con i messaggi
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(messages)
}
