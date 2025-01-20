package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) manageReactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageReaction: invalid message ID")
		return
	}

	var body struct {
		Emoji string `json:"emoji"`
		Add   bool   `json:"add"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageReaction: error decoding request body")
		return
	}

	// Aggiorna le reazioni nel DB
	err = rt.db.UpdateMessageReaction(messageId, body.Emoji, body.Add)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("manageReaction: error updating reactions")
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("manageReaction: reaction updated successfully")
}
