package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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

	// Controlla il tipo di contenuto della richiesta
	switch r.Header.Get("Content-Type") {
	case "application/json":
		// Gestisci i messaggi di testo
		var body struct {
			MessageContent string `json:"message_content"`
			Timestamp      string `json:"timestamp,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MessageContent == "" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid request body")
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
				ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid timestamp format")
				return
			}
		}

		// Crea il messaggio di testo nel gruppo
		messageId, err := rt.db.CreateGroupMessage(groupId, senderId, body.MessageContent, msgTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error saving text message to database")
			return
		}

		// Rispondi con successo
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  msgTime,
		})

	case "multipart/form-data":
		// Gestisci i messaggi con foto
		// Estrai il file immagine dalla richiesta multipart
		file, _, err := r.FormFile("photo")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error retrieving file")
			return
		}
		defer file.Close()

		// Leggi i dati del file immagine direttamente nel database come BLOB
		imageData, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error reading image file data")
			return
		}

		// Decodifica il corpo della richiesta per il messaggio
		var body struct {
			Timestamp string `json:"timestamp,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid request body")
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
				ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid timestamp format")
				return
			}
		}

		// Crea il messaggio con l'immagine nel gruppo
		messageId, err := rt.db.CreateGroupImageMessage(groupId, senderId, imageData, msgTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("sendMessageToGroup: error saving image message to database")
			return
		}

		// Rispondi con successo
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message_id": messageId,
			"status":     "sent",
			"timestamp":  msgTime,
		})

	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("sendMessageToGroup: unsupported content type")
	}
}
