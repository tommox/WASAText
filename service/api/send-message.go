package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
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
		Recipient_id   int    `json:"recipient_id"`
		MessageContent string `json:"message_content"`
		Timestamp      string `json:"timestamp,omitempty"` // Timestamp opzionale
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MessageContent == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessage: invalid request body")
		return
	}

	// Controlla se il destinatario esiste
	_, err = rt.db.CheckUserId(database.User{User_id: body.Recipient_id})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessage: recipient does not exist")
		return
	}

	// Convertire il timestamp, o usare quello corrente
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

	// Salva il messaggio
	messageId, err := rt.db.CreateMessage(senderId, body.Recipient_id, body.MessageContent, msgTime)
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
		"timestamp":  msgTime, // Aggiunto nella risposta
	})
}
