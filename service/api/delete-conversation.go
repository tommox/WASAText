package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) deleteConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationId, err := strconv.Atoi(ps.ByName("Conversation_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteConversation: invalid conversation ID")
		return
	}

	// Estrai lo User_id dal token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("deleteConversation: unauthorized user")
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deleteConversation: invalid user ID")
		return
	}

	// Verifica i permessi
	hasAccess, isGroup, err := rt.db.CheckConversationAccess(userId, conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("deleteConversation: error checking permissions")
		return
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user not allowed to delete conversation")).Error("deleteConversation: permission denied")
		return
	}

	// Se la conversazione è di gruppo, non permettiamo l'eliminazione diretta
	if isGroup {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("deleteConversation: cannot delete a group conversation directly")
		return
	}

	// Elimina la conversazione dal database
	err = rt.db.DeleteConversation(conversationId)
	if err != nil {
		if errors.Is(err, database.ErrConversationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deleteConversation: conversation not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("deleteConversation: error deleting conversation")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("deleteConversation: successfully deleted conversation ID %d", conversationId)
}
