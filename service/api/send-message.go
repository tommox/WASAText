package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) sendMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Recupera user_id dal token (mittente)
	senderIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessage: no valid token")
		return
	}
	senderId, err := strconv.Atoi(senderIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessage: converting senderId")
		return
	}

	// Legge i dati dal body (JSON)
	var body struct {
		Recipient_id   int    `json:"Recipient_id"`
		MessageContent string `json:"messageContent"`
		Timestamp      string `json:"timestamp,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validare lunghezza e pattern di `MessageContent`.
	if !validMessage(body.MessageContent) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid message format or length",
		})
		return
	}

	recipient := database.User{
		User_id: body.Recipient_id,
	}

	// Controlla se il destinatario esiste nel database
	_, err = rt.db.CheckUser(recipient)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Destinatario non trovato
			ctx.Logger.WithError(err).Error("sendMessage: recipient does not exist")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Recipient does not exist",
			})
			return
		}
		// Altro errore
		ctx.Logger.WithError(err).Error("sendMessage: error checking recipient existence")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convertire il timestamp, o se manca imposti "now"
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

	// Salva il messaggio nel DB
	messageId, err := rt.db.CreateMessage(senderId, body.Recipient_id, body.MessageContent, msgTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("sendMessage: error creating message")
		return
	}

	// Rispondi con 201, e il messageId e status
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"Message_id": messageId,
		"status":     "sent",
	})
}
