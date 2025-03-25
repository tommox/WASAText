package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
	"github.com/tommox/WASAText/service/database"
)

func (rt *_router) getGroupPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	groupId, err := strconv.Atoi(ps.ByName("Group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getGroupPhoto: invalid Group_id")
		return
	}

	// Recupera la foto del gruppo dal database
	photoData, err := rt.db.GetGroupPhoto(groupId)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("getGroupPhoto: group not found")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("getGroupPhoto: error retrieving group photo")
		}
		return
	}

	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	w.Write(photoData)
}
