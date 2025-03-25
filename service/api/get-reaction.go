package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

// manageGetReactionsHandler gestisce le richieste GET per ottenere le reazioni di un messaggio.
func (rt *_router) getReactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("Message_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("manageGetReactionsHandler: invalid Message_id")
		return
	}

	isGroupStr := r.URL.Query().Get("isGroup")
	isGroup := isGroupStr == "true"

	reactions, err := rt.db.GetReactionsForMessage(messageId, isGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("manageGetReactionsHandler: error retrieving reactions")
		return
	}
	fmt.Println(reactions)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reactions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("manageGetReactionsHandler: error encoding response")
	}
}
