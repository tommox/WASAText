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

	// Imposta il Content-Type corretto (es. image/jpeg o image/png)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(imageBytes)
}
