package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("getMessage: invalid message ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Recupera l'ID dell'utente richiedente
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("getMessage: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.Logger.WithError(err).Error("getMessage: invalid userId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verifica permessi
	hasPermission, _, err := rt.db.CheckConversationAccess(userId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getMessage: error checking permissions")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(errors.New("user not allowed to access message")).Error("getMessage: permission denied")
		return
	}

	// Recupera il messaggio dal database
	dbMsg, err := rt.db.GetMessage(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ctx.Logger.WithError(err).Error("getMessage: error retrieving message")
		return
	}

	// Converte in formato API
	apiMsg := toDatabaseMessage(dbMsg)

	// Invia la risposta al client
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(apiMsg)
}
