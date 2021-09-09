package context

import (
	"context"
	"time"
)

type detached struct {
	ctx context.Context
}

func (detached) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (detached) Done() <-chan struct{} {
	return nil
}

func (detached) Err() error {
	return nil
}

func (d detached) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

// Detach creates a detached context (without cancel and deadline) from a parent ctx.
// Example:
//	import myapp/context
//	func home(w http.ResponseWriter, r *http.Request) {
//		log.Println("home")
//		// Create a Detached context
//		detachedCtx := context.Detach(r.Context())
//		detachedDeadline, ok := detachedCtx.Deadline()
//		fmt.Println("DETACHED ---")
//		fmt.Println("Deadline():", detachedDeadline, ok) // 0001-01-01 00:00:00 +0000 UTC false
//		fmt.Println("Err():", r.Context().Err()) // <empty>
// 		backgroundTask.Run(detachedCtx) // prevents task stop when r.Context has cancelled
//	}
func Detach(ctx context.Context) context.Context {
	return detached{ctx: ctx}
}
