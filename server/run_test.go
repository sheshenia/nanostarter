package server

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

// TestRunContext test http server graceful shutdown
// it should return correct http.ErrServerClosed error
// also check if it returns error, because possible scenario it could not react on ctx.Done
func TestRunContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- Run(ctx)
	}()

	cT := time.NewTimer(time.Second * 1)
	defer cT.Stop()

	select {
	case err := <-errCh:
		// it's ok that server returned something, but need to control if it's correct error
		if err == nil || !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("context withTimeout expect error: %v, got: %v", http.ErrServerClosed, err)
		}
	// run should return error in 500 milliseconds
	// this timer fires in 1 sec, to be sure that run doesn't react on ctx.Done
	case <-cT.C:
		t.Errorf("run function doesn't react on ctx.Done")
	}

}

// go test -run TestRunContext
