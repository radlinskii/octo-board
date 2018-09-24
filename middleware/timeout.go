package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware returns 408 status code if server exceeds timeout duration to send response.
type TimeoutMiddleware struct {
	Next http.Handler
}

func (tm TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	r.WithContext(ctx)
	ch := make(chan struct{})
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
