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

func (rt *_router) sendMessageToGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid Group_id")
		return
	}

	// Estrai l'utente corrente dal Bearer Token
	senderIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: unauthorized user")
		return
	}

	senderId, err := strconv.Atoi(senderIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid sender ID")
		return
	}

	// Verifica che l'utente sia un membro del gruppo
	isMember, err := rt.db.IsGroupMember(groupId, senderId)
	if err != nil || !isMember {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: user not authorized to send message to group")
		return
	}

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		// ✉️ Messaggio testuale
		var body struct {
			MessageContent string `json:"message_content"`
			Timestamp      string `json:"timestamp,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MessageContent == "" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid request body")
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
				ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid timestamp format")
				return
			}
		}

		// Salva messaggio testuale
		messageId, err := rt.db.CreateGroupMessage(groupId, senderId, body.MessageContent, msgTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error saving text message to database")
			return
		}

		// Successo
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  msgTime,
		}); err != nil {
			ctx.Logger.WithError(err).Error("sendMessageToGroup: errore encoding JSON (text)")
		}

	} else if strings.HasPrefix(contentType, "multipart/form-data") {

		file, _, err := r.FormFile("photo")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error retrieving file")
			return
		}
		defer file.Close()

		imageData, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error reading image file data")
			return
		}

		// Salva immagine
		messageId, err := rt.db.CreateGroupImageMessage(groupId, senderId, imageData, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error saving image message to database")
			return
		}

		// Successo
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  time.Now(),
		}); err != nil {
			ctx.Logger.WithError(err).Error("sendMessageToGroup: errore encoding JSON (image)")
		}

	} else {
		// ❌ Tipo non supportato
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("sendMessageToGroup: unsupported content type")
	}
}
