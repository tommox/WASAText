package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getUsersOfGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	groupIDStr := ps.ByName("Group_id")
	ctx.Logger.Infof("Richiesta ricevuta per ottenere gli utenti del gruppo %s da: %s", groupIDStr, r.RemoteAddr)

	// Conversione dell'ID da stringa a intero
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		ctx.Logger.WithError(err).Error("ID gruppo non valido")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Recupero dei membri del gruppo dal database
	users, err := rt.db.GetGroupMembers(groupID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Errore durante il recupero dei membri del gruppo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Risposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("getUsersOfGroupHandler: errore durante l'encoding JSON")
	}
}
