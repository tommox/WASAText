package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	// Decodifica il corpo della richiesta
	var body struct {
		MessageContent string `json:"message_content"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.MessageContent == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: invalid request body")
		return
	}

	// Verifica che l'utente sia un membro del gruppo
	isMember, err := rt.db.IsGroupMember(groupId, senderId)
	if err != nil || !isMember {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: user not authorized to send message to group")
		return
	}

	// Salva il messaggio nel database
	messageId, err := rt.db.CreateGroupMessage(groupId, senderId, body.MessageContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("sendMessageToGroup: error saving message to database")
		return
	}

	// Rispondi con successo
	response := map[string]interface{}{
		"message_id": messageId,
		"status":     "sent",
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
