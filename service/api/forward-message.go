package api

import (
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

	// Recupera il messaggio per ottenere il conversation_id
	msg, err := rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("forwardMessage: failed to retrieve message")
		return
	}

	// Verifica se l'utente può accedere alla conversazione
	hasAccess, err := rt.db.CheckPrivateConversationAccess(userId, msg.Conversation_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("forwardMessage: error checking conversation access")
		return
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user has no access")).Error("forwardMessage: permission denied")
		return
	}

	// Decodifica il corpo della richiesta per ottenere la conversazione di destinazione
	var body struct {
		ConversationId int `json:"conversation_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: error decoding request body")
		return
	}

	// Inoltra il messaggio alla nuova conversazione
	newMessageId, err := rt.db.CreateMessage(userId, body.ConversationId, msg.MessageContent, time.Now())
	if err != nil {
		ctx.Logger.WithError(err).Error("forwardMessage: error creating forwarded message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Risposta con successo
	response := map[string]interface{}{
		"message_id": newMessageId,
		"status":     "sent",
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
