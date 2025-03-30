package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) checkOrCreateConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai `userId` dal token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("checkOrCreateConversation: no valid token")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("checkOrCreateConversation: invalid user ID")
		return
	}

	// Leggi il `recipient_id` dal corpo della richiesta JSON
	var requestData struct {
		RecipientId int       `json:"recipient_id"`
		MessageId   int       `json:"message_id"`
		Timestamp   time.Time `json:"timestamp"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("checkOrCreateConversation: invalid JSON")
		return
	}

	// Verifica se esiste già una conversazione
	conversationId, err := rt.db.CheckExistingConversation(userId, requestData.RecipientId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("checkOrCreateConversation: error checking existing conversation")
		return
	}

	if conversationId == 0 {
		// Se non esiste, crea la nuova conversazione
		conversationId, err = rt.db.UpdateOrCreateConversation(userId, requestData.RecipientId, 0, time.Now(), false, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("checkOrCreateConversation: error creating conversation")
			return
		}
	} else if requestData.MessageId != 0 {
		// Se invece la conversazione esiste, aggiorna l'ultimo messaggio
		conversationId, err = rt.db.UpdateOrCreateConversation(userId, requestData.RecipientId, requestData.MessageId, requestData.Timestamp, false, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("checkOrCreateConversation: error updating conversation")
			return
		}
	}

	// Risposta con l'ID della conversazione (sia esistente che appena creata)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conversation_id": conversationId,
	})
}
