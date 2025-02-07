package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) setMyPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	user_id, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid User_id")
		return
	}

	requestingUser_id_str, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("Unauthorized access")
		return
	}

	requestingUser_id, err := strconv.Atoi(requestingUser_id_str)
	if err != nil || user_id != requestingUser_id {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("Unauthorized: user_id mismatch")
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

	// Leggi i byte del file
	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error reading file data")
		return
	}

	// Salva l'immagine nel database
	err = rt.db.UpdateUserPhoto(user_id, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Error updating photo in database")
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("Photo updated successfully for user_id: %d", user_id)
}
