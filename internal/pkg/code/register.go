package code

import "net/http"

func init() {
	register(ErrSuccess, http.StatusOK, "success")
	register(ErrBind, http.StatusBadRequest, "invalid request")
}