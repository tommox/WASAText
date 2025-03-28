package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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

	messageType := r.URL.Query().Get("type")
	if messageType != "private" && messageType != "group" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(errors.New("invalid message type")).Error("getMessage: invalid type parameter")
		return
	}

	if messageType == "private" {
		// Verifica accesso
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

		// Recupera messaggio
		dbMsg, err := rt.db.GetMessage(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("getMessage: message not found")
			return
		}

		// Se immagine, recupera anche imageData + timestamp
		if dbMsg.ImageData != nil {
			imageData, timestamp, err := rt.db.GetMessageImage(messageId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("getMessage: error retrieving image")
				return
			}
			response := map[string]interface{}{
				"image_data": base64.StdEncoding.EncodeToString(imageData),
				"timestamp":  timestamp,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		// Altrimenti, messaggio di testo
		w.Header().Set("Content-Type", "application/json")
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

		groupMsg, err := rt.db.GetGroupMessage(groupConv.Group_id, messageId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("getMessage: error retrieving group message")
			return
		}

		if groupMsg.ImageData != nil {
			imageData, timestamp, err := rt.db.GetGroupMessageImage(messageId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("getMessage: error retrieving group image")
				return
			}
			response := map[string]interface{}{
				"image_data": base64.StdEncoding.EncodeToString(imageData),
				"timestamp":  timestamp,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(groupMsg)
		return
	}
}
