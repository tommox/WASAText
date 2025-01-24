package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getConversationsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Recupera messaggi normali
	messages, err := rt.db.GetAllMessages()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getConversations: error retrieving normal messages")
		return
	}

	// Recupera messaggi di gruppo
	groupMessages, err := rt.db.GetAllGroupMessages()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getConversations: error retrieving group messages")
		return
	}

	// Converte i messaggi in formato API
	var apiMessages []Message
	for _, msg := range messages {
		apiMessages = append(apiMessages, Message{
			Message_id:     msg.Message_id,
			Sender_id:      msg.Sender_id,
			Recipient_id:   msg.Recipient_id,
			MessageContent: msg.MessageContent,
			Timestamp:      msg.Timestamp,
		})
	}

	var apiGroupMessages []GroupMessage
	for _, msg := range groupMessages {
		apiGroupMessages = append(apiGroupMessages, GroupMessage{
			GroupMessage_id: msg.GroupMessage_id,
			Group_id:        msg.Group_id,
			Sender_id:       msg.Sender_id,
			MessageContent:  msg.MessageContent,
			Timestamp:       msg.Timestamp,
		})
	}

	// Costruisci la risposta
	response := map[string]interface{}{
		"private_messages": apiMessages,
		"group_messages":   apiGroupMessages,
	}

	// Rispondi
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
