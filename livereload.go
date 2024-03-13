// Package livereload provides some JavaScript and corresponding http.Handler
// to enable full page reloads when the http.Server restarts.
package livereload

import (
	"fmt"
	"net/http"
	"time"
)

var now = time.Now()

// JS should included on the webpages that need livereload. Including it will
// automatically setup the necessary functionality. Make sure to setup the
// corresponding Handler.
const JS = `
(() => {
  let token;
  let att = 0;
  let pending;
  const schedule = () => {
    clearTimeout(pending);
    pending = setTimeout(monitor, Math.min(50 * Math.pow(1.05, att++), 5000));
  };
  const monitor = () => {
    const es = new EventSource('/livereload');
    window.addEventListener('beforeunload', () => es.close());
    es.addEventListener('token', ev => {
      if (!token) {
        token = event.data;
      } else if (token !== event.data) {
        es.close();
        location.reload();
      }
    })
    es.addEventListener('error', ev => {
      es.close();
      schedule();
    })
  };
  monitor();
})();
`

// Handler should be mounted to the path "/livereload". It will be used by
// the corresponding JS.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	shutdown := make(chan struct{})
	r.Context().Value(http.ServerContextKey).(*http.Server).RegisterOnShutdown(func() {
		close(shutdown)
	})
	done := r.Context().Done()

	rc := http.NewResponseController(w)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		_, err := fmt.Fprintf(w, "event: token\ndata: %d\n\n", now.UnixMicro())
		if err != nil {
			return
		}
		if err := rc.Flush(); err != nil {
			return
		}
		select {
		case <-shutdown:
			return
		case <-done:
			return
		case <-ticker.C:
			continue
		}
	}
}
