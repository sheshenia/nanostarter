package server

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// TestSimpleCommandHandler tests simple command handler
// send command to "list" this folder's content, it should contain current test file name
func TestSimpleCommandHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/command?cmd="+url.QueryEscape("ls -lh"), nil)
	w := httptest.NewRecorder()
	simpleCommandHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if !bytes.Contains(data, []byte("run_test.go")) {
		t.Error("expected to find run_test.go in folder content")
	}
}
