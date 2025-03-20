package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("getMessage: invalid message ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("getMessage: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.Logger.WithError(err).Error("getMessage: invalid userId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 🏷️ Controllo che il tipo sia valido
	messageType := r.URL.Query().Get("type")
	if messageType != "private" && messageType != "group" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(errors.New("invalid message type")).Error("getMessage: invalid type parameter")
		return
	}

	if messageType == "private" {
		conversationId, err := rt.db.GetConversationIdByMessageId(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("getMessage: conversation not found")
			return
		}

		hasPermission, err := rt.db.CheckPrivateConversationAccess(userId, conversationId)
		if err != nil || !hasPermission {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(err).Error("getMessage: no access to private chat")
			return
		}

		dbMsg, err := rt.db.GetMessage(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("getMessage: message not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(toDatabaseMessage(dbMsg))
		return
	}

	if messageType == "group" {
		groupConv, err := rt.db.GetGroupByMessageId(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("getMessage: group not found")
			return
		}

		isMember, err := rt.db.IsGroupMember(groupConv.Group_id, userId)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(err).Error("getMessage: no access to group chat")
			return
		}

		lastGroupMessage, err := rt.db.GetGroupMessage(groupConv.Group_id, messageId)
		fmt.Println("message:", lastGroupMessage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("getMessage: error retrieving group messages")
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(lastGroupMessage)
		return
	}
}
