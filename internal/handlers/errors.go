package handlers

import "net/http"

func (h *UserHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"errors": message}

	err := h.writeJSON(w, status, env, nil)
	if err != nil {
		h.Logger.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Print(err)
	message := "the server encountered internal error"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}
