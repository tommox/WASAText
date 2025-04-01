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
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: invalid messageId")
		return
	}

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

	messageType := r.URL.Query().Get("type")
	if messageType != messageTypePrivate && messageType != messageTypeGroup {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(errors.New("invalid message type")).Error("forwardMessage: invalid type")
		return
	}

	var body struct {
		ConversationId *int `json:"conversation_id,omitempty"`
		GroupId        *int `json:"group_id,omitempty"`
		IsForward      bool `json:"isForward,omitempty"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("forwardMessage: error decoding body")
		return
	}

	// ✅ Origine: Messaggio privato
	if messageType == messageTypePrivate {
		msg, err := rt.db.GetMessage(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("forwardMessage: private message not found")
			return
		}

		hasAccess, err := rt.db.CheckPrivateConversationAccess(userId, msg.Conversation_id)
		if err != nil || !hasAccess {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(err).Error("forwardMessage: no access to source private conversation")
			return
		}

		if body.ConversationId != nil {
			// ✅ Destinazione: conversazione privata
			var newMessageId int
			if msg.ImageData != nil {
				newMessageId, err = rt.db.CreateImageMessage(userId, *body.ConversationId, msg.ImageData, time.Now())
			} else {
				newMessageId, err = rt.db.CreateMessage(userId, *body.ConversationId, msg.MessageContent, time.Now(), nil, body.IsForward)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("forwardMessage: failed to forward private → private")
				return
			}
			if body.IsForward {
				err := rt.db.MarkIsForward(newMessageId, body.IsForward)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("forwardMessage: failed to mark message as forwarded")
					return
				}
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message_id": newMessageId,
				"status":     "sent",
			}); err != nil {
				ctx.Logger.WithError(err).Error("forwardMessage: errore nell'encoding JSON (private → private)")
			}
			return
		} else if body.GroupId != nil {
			// ✅ Destinazione: gruppo
			isTargetMember, err := rt.db.IsGroupMember(*body.GroupId, userId)
			if err != nil || !isTargetMember {
				w.WriteHeader(http.StatusForbidden)
				ctx.Logger.WithError(err).Error("forwardMessage: no access to target group (private → group)")
				return
			}
			// Inoltra il messaggio al gruppo senza impostare l'IsReply
			var newMessageId int
			if msg.ImageData != nil {
				newMessageId, err = rt.db.CreateGroupImageMessage(*body.GroupId, userId, msg.ImageData, time.Now())
			} else {
				newMessageId, err = rt.db.CreateGroupMessage(*body.GroupId, userId, msg.MessageContent, time.Now(), nil, body.IsForward)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("forwardMessage: failed to forward private → group")
				return
			}
			if body.IsForward {
				err := rt.db.MarkIsForwardGroup(newMessageId, body.IsForward)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("forwardMessage: failed to mark message as forwarded for group")
					return
				}
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message_id": newMessageId,
				"status":     "sent",
			}); err != nil {
				ctx.Logger.WithError(err).Error("forwardMessage: errore nell'encoding JSON (private → group)")
			}
			return
		}
	}

	// ✅ Origine: Messaggio di gruppo
	if messageType == messageTypeGroup {
		group, err := rt.db.GetGroupByMessageId(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("forwardMessage: group not found for message")
			return
		}

		isMember, err := rt.db.IsGroupMember(group.Group_id, userId)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(errors.New("user not in source group")).Error("forwardMessage: source group access denied")
			return
		}

		msg, err := rt.db.GetGroupMessage(group.Group_id, messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("forwardMessage: original group message not found")
			return
		}

		if body.ConversationId != nil {
			// ✅ Destinazione: conversazione privata
			hasAccess, err := rt.db.CheckPrivateConversationAccess(userId, *body.ConversationId)
			if err != nil || !hasAccess {
				w.WriteHeader(http.StatusForbidden)
				ctx.Logger.WithError(err).Error("forwardMessage: no access to target private (group → private)")
				return
			}

			var newMessageId int
			if msg.ImageData != nil {
				newMessageId, err = rt.db.CreateImageMessage(userId, *body.ConversationId, msg.ImageData, time.Now())
			} else {
				newMessageId, err = rt.db.CreateMessage(userId, *body.ConversationId, msg.MessageContent, time.Now(), nil, body.IsForward)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("forwardMessage: failed to forward group → private")
				return
			}
			if body.IsForward {
				err := rt.db.MarkIsForward(newMessageId, body.IsForward)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("forwardMessage: failed to mark message as forwarded for private")
					return
				}
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message_id": newMessageId,
				"status":     "sent",
			}); err != nil {
				ctx.Logger.WithError(err).Error("forwardMessage: errore nell'encoding JSON (group → private)")
			}
			return
		} else if body.GroupId != nil {
			// ✅ Destinazione: altro gruppo
			isTargetMember, err := rt.db.IsGroupMember(*body.GroupId, userId)
			if err != nil || !isTargetMember {
				w.WriteHeader(http.StatusForbidden)
				ctx.Logger.WithError(err).Error("forwardMessage: no access to target group")
				return
			}

			var newMessageId int
			if msg.ImageData != nil {
				newMessageId, err = rt.db.CreateGroupImageMessage(*body.GroupId, userId, msg.ImageData, time.Now())
			} else {
				newMessageId, err = rt.db.CreateGroupMessage(*body.GroupId, userId, msg.MessageContent, time.Now(), nil, body.IsForward)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("forwardMessage: failed to forward group → group")
				return
			}
			if body.IsForward {
				err := rt.db.MarkIsForwardGroup(newMessageId, body.IsForward)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("forwardMessage: failed to mark message as forwarded for group")
					return
				}
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message_id": newMessageId,
				"status":     "sent",
			}); err != nil {
				ctx.Logger.WithError(err).Error("forwardMessage: errore nell'encoding JSON (group → group)")
			}
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	ctx.Logger.Error("forwardMessage: missing or invalid destination")
}
