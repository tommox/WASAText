package api

import (
	"errors"
	"fmt"
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

	fmt.Println(userId)
	// Verifica i permessi
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

	// Elimina il messaggio dal database
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

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("deleteMessage: successfully deleted message ID %d", messageId)
}
