package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) addReactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addReaction: invalid Message_id")
		return
	}

	// Estrai l'utente corrente dal Bearer Token
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("addReaction: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addReaction: invalid userId")
		return
	}

	var body struct {
		Emoji   string `json:"emoji"`
		IsGroup bool   `json:"isGroup"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Emoji == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("addReaction: invalid request body")
		return
	}
	err = rt.db.AddReaction(messageId, userId, body.Emoji, body.IsGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("addReaction: error adding/updating reaction")
		return
	}
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("addReaction: reaction successfully added for user %d on message %d", userId, messageId)
}
