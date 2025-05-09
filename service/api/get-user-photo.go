package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) getUserPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	user_id, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid User_id")
		return
	}

	// Recupera la foto dal database
	imageBytes, err := rt.db.GetUserPhoto(user_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.WithError(err).Error("Photo not found")
		return
	}

	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(imageBytes); err != nil {
		ctx.Logger.WithError(err).Error("Error writing image to response")
	}
}
