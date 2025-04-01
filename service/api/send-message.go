package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
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

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		// ✉️ Messaggio testuale
		var body struct {
			ConversationId int    `json:"conversation_id"`
			MessageContent string `json:"message_content"`
			Timestamp      string `json:"timestamp,omitempty"`
			IsReply        *int   `json:"isReply,omitempty"` // Aggiunta gestione per la risposta
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MessageContent == "" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessage: invalid request body")
			return
		}

		// Verifica accesso alla conversazione
		hasAccess, err := rt.db.CheckPrivateConversationAccess(senderId, body.ConversationId)
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

		// Timestamp
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

		// Salva messaggio
		messageId, err := rt.db.CreateMessage(senderId, body.ConversationId, body.MessageContent, msgTime, body.IsReply, false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessage: failed to create message")
			return
		}

		if body.IsReply != nil {
			err := rt.db.MarkIsReply(messageId, *body.IsReply)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("sendMessage: failed to update IsReply")
				return
			}
		}

		// Successo
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  msgTime,
		}); err != nil {
			ctx.Logger.WithError(err).Error("sendMessage: errore encoding JSON (text)")
		}

	} else if strings.HasPrefix(contentType, "multipart/form-data") {

		file, _, err := r.FormFile("photo")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessage: error retrieving file")
			return
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessage: error reading image file data")
			return
		}

		// Estrai campi form
		conversationIdStr := r.FormValue("conversation_id")

		conversationId, err := strconv.Atoi(conversationIdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessage: invalid conversation_id")
			return
		}

		// Verifica accesso alla conversazione
		hasAccess, err := rt.db.CheckPrivateConversationAccess(senderId, conversationId)
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

		// Salva immagine
		messageId, err := rt.db.CreateImageMessage(senderId, conversationId, imageData, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessage: failed to create image message")
			return
		}

		isReplyStr := r.FormValue("isReply")
		if isReplyStr != "" {
			if replyId, err := strconv.Atoi(isReplyStr); err == nil {
				err := rt.db.MarkIsReply(messageId, replyId)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("sendMessage: failed to update IsReply for image")
					return
				}
			}
		}

		// Successo
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  time.Now(),
		}); err != nil {
			ctx.Logger.WithError(err).Error("sendMessage: errore encoding JSON (image)")
		}

	} else {
		// ❌ Tipo non supportato
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("sendMessage: unsupported content type")
	}
}
