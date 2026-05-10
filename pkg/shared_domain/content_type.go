package shared_domain

import "context"

type ResponseFormat int

const (
	ResponseFormatHTML ResponseFormat = iota
	ResponseFormatJSON
	ResponseFormatFullPage
)

type contextKey string

const responseFormatKey contextKey = "response_format"

func WithResponseFormat(ctx context.Context, format ResponseFormat) context.Context {
	return context.WithValue(ctx, responseFormatKey, format)
}

func GetResponseFormat(ctx context.Context) ResponseFormat {
	v, ok := ctx.Value(responseFormatKey).(ResponseFormat)
	if !ok {
		return ResponseFormatFullPage
	}
	return v
}
