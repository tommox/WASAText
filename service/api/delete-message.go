package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) deleteMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteMessage: invalid message ID")
		return
	}

	// Estrai lo User_id dal token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("deleteMessage: unauthorized user")
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteMessage: invalid user ID")
		return
	}

	// Controllo parametro tipo
	messageType := r.URL.Query().Get("type")
	if messageType != "private" && messageType != "group" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(errors.New("invalid message type")).Error("deleteMessage: invalid type parameter")
		return
	}

	if messageType == "private" {
		// 🔐 Controllo permessi per messaggi privati
		hasPermission, err := rt.db.CheckUserPermission(userId, messageId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("deleteMessage: error checking permissions")
			return
		}
		if !hasPermission {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(errors.New("user not allowed to delete message")).Error("deleteMessage: permission denied")
			return
		}

		// 🗑️ Elimina il messaggio privato
		err = rt.db.DeleteMessage(messageId)
		if err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				w.WriteHeader(http.StatusNotFound)
				ctx.Logger.WithError(err).Error("deleteMessage: message not found")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("deleteMessage: error deleting message")
			}
			return
		}

		ctx.Logger.Infof("deleteMessage: successfully deleted private message ID %d", messageId)
		w.WriteHeader(http.StatusOK)
		return
	}

	if messageType == "group" {
		groupConv, err := rt.db.GetGroupByMessageId(messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deleteMessage: group not found")
			return
		}

		// ✅ Solo l'autore o un admin del gruppo può cancellare
		isMember, err := rt.db.IsGroupMember(groupConv.Group_id, userId)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(err).Error("deleteMessage: user not member of group")
			return
		}

		// Ottieni il messaggio di gruppo
		groupMsg, err := rt.db.GetGroupMessage(groupConv.Group_id, messageId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deleteMessage: group message not found")
			return
		}

		// Verifica se è l’autore oppure admin
		isAdmin, err := rt.db.IsGroupAdmin(groupConv.Group_id, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("deleteMessage: error checking admin")
			return
		}

		if userId != groupMsg.Sender_id && !isAdmin {
			w.WriteHeader(http.StatusForbidden)
			ctx.Logger.WithError(errors.New("user not allowed to delete group message")).Error("deleteMessage: permission denied")
			return
		}

		// Elimina il messaggio di gruppo
		err = rt.db.DeleteGroupMessage(messageId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("deleteMessage: error deleting group message")
			return
		}

		ctx.Logger.Infof("deleteMessage: successfully deleted group message ID %d", messageId)
		w.WriteHeader(http.StatusOK)
		return
	}
}
