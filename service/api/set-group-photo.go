package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) setGroupPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il Group_id dal percorso
	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setGroupPhoto: invalid Group_id")
		return
	}

	// Estrai l'utente loggato dal Bearer Token
	adminIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("setGroupPhoto: unauthorized user")
		return
	}

	adminId, err := strconv.Atoi(adminIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setGroupPhoto: invalid admin ID")
		return
	}

	// Controlla che l'utente sia admin del gruppo
	isAdmin, err := rt.db.IsGroupAdmin(groupId, adminId)
	if err != nil || !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("setGroupPhoto: user is not an admin")
		return
	}

	// Leggi la foto dal body
	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setGroupPhoto: invalid photo data")
		return
	}

	// Aggiorna la foto del gruppo nel database
	err = rt.db.UpdateGroupPhoto(groupId, data)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("setGroupPhoto: group not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("setGroupPhoto: error updating group photo")
		}
		return
	}

	// Risposta di successo
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("setGroupPhoto: photo updated successfully for group_id %d", groupId)
}
