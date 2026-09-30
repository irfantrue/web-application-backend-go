package ctxval

import (
	"context"
	"time"
)

type contextKey string

const timezoneKey contextKey = "user_timezone"

func GetTimezone(ctx context.Context) *time.Location {
	if loc, ok := ctx.Value(timezoneKey).(*time.Location); ok {
		return loc
	}
	return time.UTC
}

func SetTimezone(ctx context.Context, loc *time.Location) context.Context {
	return context.WithValue(ctx, timezoneKey, loc)
}
