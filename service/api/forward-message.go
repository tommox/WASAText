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
)

func (rt *_router) forwardMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il messageId dalla richiesta
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: invalid messageId")
		return
	}

	// Estrai l'utente corrente dal Bearer Token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("forwardMessage: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: invalid userId")
		return
	}

	// Recupera il messaggio originale
	originalMessage, err := rt.db.GetMessage(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("forwardMessage: message not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("forwardMessage: error retrieving message")
		}
		return
	}

	// Verifica se l'utente ha i permessi per accedere al messaggio
	if originalMessage.Recipient_id != userId && originalMessage.Sender_id != userId {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user not allowed to access message")).Error("forwardMessage: permission denied")
		return
	}

	// Decodifica il corpo della richiesta
	var body struct {
		Recipient_id int `json:"Recipient_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: error decoding request body")
		return
	}

	// Salva il messaggio inoltrato nel DB
	newMessageId, err := rt.db.CreateMessage(userId, body.Recipient_id, originalMessage.MessageContent, time.Now())
	if err != nil {
		ctx.Logger.WithError(err).Error("forwardMessage: error creating forwarded message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Risposta con successo
	response := map[string]interface{}{
		"Message_id": newMessageId,
		"status":     "sent",
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
