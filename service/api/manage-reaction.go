package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) manageReactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageReaction: invalid Message_id")
		return
	}

	// Estrai l'utente corrente dal Bearer Token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("manageReaction: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageReaction: invalid userId")
		return
	}

	var body struct {
		Emoji   string `json:"emoji"`
		Add     bool   `json:"add"`
		IsGroup bool   `json:"isGroup"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Emoji == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageReaction: invalid request body")
		return
	}

	if body.Add {
		err = rt.db.AddReaction(messageId, userId, body.Emoji, body.IsGroup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("manageReaction: error adding/updating reaction")
			return
		}
	} else {
		err = rt.db.RemoveReaction(messageId, userId, body.IsGroup)
		if err != nil {
			if errors.Is(err, database.ErrReactionNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			ctx.Logger.WithError(err).Error("manageReaction: error removing reaction")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("manageReaction: reaction successfully managed for user %d on message %d", userId, messageId)
}
