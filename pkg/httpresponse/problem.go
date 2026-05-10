package httpresponse

import "net/http"

const ContentTypeProblemJSON = "application/problem+json"

type ProblemDetail struct {
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Status     int            `json:"status"`
	Detail     string         `json:"detail,omitempty"`
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"extensions,omitempty"`
}

func NewProblemDetail(problemType, title string, status int) ProblemDetail {
	return ProblemDetail{
		Type:   problemType,
		Title:  title,
		Status: status,
	}
}

func (p ProblemDetail) WithDetail(detail string) ProblemDetail {
	p.Detail = detail
	return p
}

func (p ProblemDetail) WithInstance(instance string) ProblemDetail {
	p.Instance = instance
	return p
}

func (p ProblemDetail) WithExtension(key string, value any) ProblemDetail {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}
	p.Extensions[key] = value
	return p
}

func (p ProblemDetail) Error() string {
	if p.Detail != "" {
		return p.Detail
	}
	return p.Title
}

const (
	TypeBadRequest          = "https://httpstatuses.com/400"
	TypeUnauthorized        = "https://httpstatuses.com/401"
	TypeForbidden           = "https://httpstatuses.com/403"
	TypeNotFound            = "https://httpstatuses.com/404"
	TypeMethodNotAllowed    = "https://httpstatuses.com/405"
	TypeUnprocessableEntity = "https://httpstatuses.com/422"
	TypeConflict            = "https://httpstatuses.com/409"
	TypeInternalServerError = "https://httpstatuses.com/500"
	TypeServiceUnavailable  = "https://httpstatuses.com/503"
)

var (
	ErrBadRequest          = NewProblemDetail(TypeBadRequest, "Bad Request", http.StatusBadRequest)
	ErrUnauthorized        = NewProblemDetail(TypeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden           = NewProblemDetail(TypeForbidden, "Forbidden", http.StatusForbidden)
	ErrNotFound            = NewProblemDetail(TypeNotFound, "Not Found", http.StatusNotFound)
	ErrMethodNotAllowed    = NewProblemDetail(TypeMethodNotAllowed, "Method Not Allowed", http.StatusMethodNotAllowed)
	ErrConflict            = NewProblemDetail(TypeConflict, "Conflict", http.StatusConflict)
	ErrUnprocessableEntity = NewProblemDetail(TypeUnprocessableEntity, "Unprocessable Entity", http.StatusUnprocessableEntity)
	ErrInternalServerError = NewProblemDetail(TypeInternalServerError, "Internal Server Error", http.StatusInternalServerError)
	ErrServiceUnavailable  = NewProblemDetail(TypeServiceUnavailable, "Service Unavailable", http.StatusServiceUnavailable)
)
