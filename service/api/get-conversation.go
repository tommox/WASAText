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
	conversationId, err := strconv.Atoi(ps.ByName("Conversation_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getConversation: invalid Conversation_id")
		return
	}

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

	conversationType := r.URL.Query().Get("type")

	if conversationType == "group" {
		isGroup, err := rt.db.CheckGroupConversationAccess(userId, conversationId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("getConversation: error checking group access")
			return
		}
		if isGroup {
			messages, err := rt.db.GetGroupConversationMessages(conversationId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("getConversation: failed to retrieve group messages")
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(messages)
			return
		}
	}

	// ⚡ Se è una conversazione privata, controlliamo le chat normali
	if conversationType == "private" {
		isPrivate, err := rt.db.CheckPrivateConversationAccess(userId, conversationId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("getConversation: error checking private access")
			return
		}
		if isPrivate {
			messages, err := rt.db.GetConversationMessages(conversationId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("getConversation: failed to retrieve private messages")
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(messages)
			return
		}
	}

	w.WriteHeader(http.StatusForbidden)
	ctx.Logger.WithError(errors.New("user has no access")).Error("getConversation: user has no access to this conversation")
}
