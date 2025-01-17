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
		return
	}

	requestingUser_id_str, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	requestingUser_id, err := strconv.Atoi(requestingUser_id_str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user_id != requestingUser_id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateUserPhoto(user_id, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
