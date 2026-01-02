package main

import (
	"chonlatee/mw/server/ctxutil"
	"chonlatee/mw/server/trace"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

func Trace(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID, err := uuid.Parse(r.Header.Get("X-Trace-Id"))
		if err != nil {
			traceID = uuid.New()
		}

		reqID, err := uuid.Parse(r.Header.Get("X-Request-Id"))
		if err != nil {
			reqID = uuid.New()
		}

		trc := trace.Trace{TraceID: traceID, RequestID: reqID}
		ctx = ctxutil.WithValue(ctx, trc)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
}

func Log(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var prefix string
		trc, ok := ctxutil.Value[trace.Trace](r.Context())
		if ok {
			prefix = fmt.Sprintf("%s %s: [%s %s]: ", r.Method, r.URL, trc.TraceID, trc.RequestID)
		} else {
			prefix = fmt.Sprintf("%s %s", r.Method, r.URL)
		}
		logger := log.New(os.Stderr, prefix, log.LstdFlags)
		ctx := ctxutil.WithValue(r.Context(), logger)
		r = r.Clone(ctx)
		h.ServeHTTP(w, r)
	}
}

func RecordResposne(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rrw := &RecordingResponseWriter{RW: w}
		start := time.Now()
		h.ServeHTTP(rrw, r)
		elapsed := time.Since(start)

		logger, ok := ctxutil.Value[*log.Logger](r.Context())
		if !ok {
			log.Printf("%s %s: %d %s: %d bytes in %s", r.Method, r.URL, rrw.StatusCode,
				http.StatusText(rrw.StatusCode), rrw.Bytes, elapsed)
			return
		}
		logger.Printf("%d %s: %d bytes in %s", rrw.StatusCode, http.StatusText(rrw.StatusCode), rrw.Bytes, elapsed)
	}
}

func Recovery(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger, ok := ctxutil.Value[*log.Logger](r.Context())
				if !ok {
					log.Printf("%s %s: panic %v\n%s", r.Method, r.URL, err, stack)
				} else {
					logger.Printf("panic: %v\n%s", err, stack)
				}

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal server error"))
			}
		}()
		h.ServeHTTP(w, r)
	}
}

type RecordingResponseWriter struct {
	RW         http.ResponseWriter
	StatusCode int
	Bytes      int // total bytes written
}

func (w *RecordingResponseWriter) WriteHeader(statusCode int) {
	if w.StatusCode == 0 {
		w.StatusCode = statusCode
	}
	w.RW.WriteHeader(statusCode)
}

func (w *RecordingResponseWriter) Header() http.Header {
	return w.RW.Header()
}

func (w *RecordingResponseWriter) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.RW.Write(b)
	w.Bytes += n
	return n, err
}

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/time":
			fmt.Fprintln(w, time.Now().Format(time.RFC3339))
		case "/panic":
			panic("oh my god!")
		default:
			http.NotFound(w, r)
		}
	}

	h = RecordResposne(h)
	// h = Recovery(h) // for response something to client if not do this it will not response anything to client
	h = Log(h)
	h = Trace(h)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		Handler:           h,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		ReadHeaderTimeout: 200 * time.Millisecond,
	}

	log.Printf("listening on %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
