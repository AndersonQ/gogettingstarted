package tracking

import (
	"context"

	"github.com/google/uuid"
)

type key struct{}

var ctxKey = key{}

func ContextWithID(ctx context.Context) context.Context {
	var id string

	// TODO: handle the error properly if you are planing to copy this!
	if uuid, err := uuid.NewUUID(); err == nil {
		id = uuid.String()
	}

	return context.WithValue(ctx, ctxKey, id)
}

func IdFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey).(string)
	return id
}
