package analytics

import "context"

type contextKey string

const postHogKey contextKey = "posthog_api_key"

func GetPostHogApiKey(ctx context.Context) string {
	if apiKey, ok := ctx.Value(postHogKey).(string); ok {
		return apiKey
	}
	return ""
}

func WithPostHogApiKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, postHogKey, apiKey)
}
