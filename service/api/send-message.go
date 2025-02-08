package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) sendMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai `Sender_id` dal token
	senderIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessage: no valid token")
		return
	}

	senderId, err := strconv.Atoi(senderIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessage: invalid sender ID")
		return
	}

	// Decodifica il corpo della richiesta
	var body struct {
		ConversationId int    `json:"conversation_id"`
		MessageContent string `json:"message_content"`
		Timestamp      string `json:"timestamp,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MessageContent == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessage: invalid request body")
		return
	}

	// Verifica l'accesso alla conversazione
	hasAccess, _, err := rt.db.CheckConversationAccess(senderId, body.ConversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("sendMessage: error checking conversation access")
		return
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessage: user has no access to this conversation")
		return
	}

	// Converti il timestamp o usa quello corrente
	var msgTime time.Time
	if body.Timestamp == "" {
		msgTime = time.Now()
	} else {
		msgTime, err = time.Parse(time.RFC3339, body.Timestamp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessage: invalid timestamp format")
			return
		}
	}

	// Crea il nuovo messaggio
	messageId, err := rt.db.CreateMessage(senderId, body.ConversationId, body.MessageContent, msgTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("sendMessage: failed to create message")
		return
	}

	// Risposta con successo
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message_id": messageId,
		"status":     "sent",
		"timestamp":  msgTime,
	})
}
