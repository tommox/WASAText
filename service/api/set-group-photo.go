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
	userIdStr, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("setGroupPhoto: unauthorized user")
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setGroupPhoto: invalid user ID")
		return
	}

	// Controlla che l'utente sia memebro del gruppo
	isMember, err := rt.db.IsGroupMember(groupId, userId)
	if err != nil || !isMember {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("setGroupPhoto: user is not an user")
		return
	}

	// Estrai il file dall'upload multipart
	file, _, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error retrieving file")
		return
	}
	defer file.Close()

	// Leggi la foto dal body
	data, err := io.ReadAll(file)
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
