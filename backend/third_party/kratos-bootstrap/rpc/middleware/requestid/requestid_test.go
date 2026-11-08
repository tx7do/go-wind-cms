package middleware

import (
	"context"
	"strings"
	"testing"
)

func TestRequestIDMiddleware_Generate(t *testing.T) {
	mw := NewRequestIDMiddleware()
	h := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		// read from context
		id := GetRequestID(ctx)
		if id == "" {
			t.Fatalf("expected request id in context")
		}
		if !strings.Contains(id, "-") {
			t.Fatalf("expected uuid-like id, got: %s", id)
		}
		return nil, nil
	})

	// call with empty ctx
	_, err := h(context.Background(), nil)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
}

func TestRequestIDMiddleware_Preserve(t *testing.T) {
	custom := "my-id-123"
	mw := NewRequestIDMiddleware()
	h := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		id := GetRequestID(ctx)
		if id != custom {
			t.Fatalf("expected preserved id %s, got %s", custom, id)
		}
		return nil, nil
	})

	ctx := context.WithValue(context.Background(), ctxKeyRequestID{}, custom)
	_, err := h(ctx, nil)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
}
