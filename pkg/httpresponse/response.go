package httpresponse

import (
	"net/http"
)

type Response struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

func WriteOK(w http.ResponseWriter, r *http.Request, data any) {
	writeJSON(w, http.StatusOK, Response{Data: data})
}

func WriteCreated(w http.ResponseWriter, r *http.Request, data any) {
	writeJSON(w, http.StatusCreated, Response{Data: data})
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteWithMeta(w http.ResponseWriter, r *http.Request, status int, data, meta any) {
	writeJSON(w, status, Response{Data: data, Meta: meta})
}

func WriteProblem(w http.ResponseWriter, problem ProblemDetail) {
	w.Header().Set("Content-Type", ContentTypeProblemJSON)
	w.WriteHeader(problem.Status)
	writeJSONBody(w, problem)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrBadRequest.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteNotFound(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrNotFound.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteUnauthorized(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrUnauthorized.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteForbidden(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrForbidden.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteConflict(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrConflict.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteUnprocessableEntity(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrUnprocessableEntity.WithDetail(detail).WithInstance(r.URL.Path))
}

func WriteInternalServerError(w http.ResponseWriter, r *http.Request, detail string) {
	WriteProblem(w, ErrInternalServerError.WithDetail(detail).WithInstance(r.URL.Path))
}
