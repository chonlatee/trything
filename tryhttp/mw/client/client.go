package main

import (
	"chonlatee/mw/ctxutil"
	"chonlatee/mw/trace"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var _ http.RoundTripper = RoundTripFunc(nil)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func logExec(name string) func() {
	log.Printf("middleware: begin: %s", name)
	return func() {
		defer log.Printf("middleware: end %s", name)
	}
}

func TimeRequest(rt http.RoundTripper) RoundTripFunc {
	return func(r *http.Request) (*http.Response, error) {
		defer logExec("TimeRequest")()
		start := time.Now()

		resp, err := rt.RoundTrip(r)
		if err != nil {
			log.Printf("%s %s: %d %s errored after %s", r.Method, r.URL, resp.StatusCode,
				http.StatusText(resp.StatusCode), time.Since(start))
			return nil, err
		}

		return resp, nil
	}
}

func RetryOn5xx(rt http.RoundTripper, wait time.Duration, tries int) RoundTripFunc {
	if tries <= 1 {
		panic("n must be > 1")
	}
	if wait <= 0 {
		panic("wait must be > 0")
	}

	return func(r *http.Request) (*http.Response, error) {
		defer logExec("RetryOn5xx")()

		var retryErrs error

		for retry := uint(0); retry < uint(tries); retry++ {
			if retry > 0 {
				log.Printf("retry attemp: #%d", retry)
				time.Sleep(wait << retry)
			}

			resp, err := rt.RoundTrip(r)
			if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
				retryErrs = errors.Join(retryErrs, err)
				continue
			}

			if retryErrs != nil {
				return nil, fmt.Errorf("failed after %d retries: %w", retry, retryErrs)
			}

			switch sc := resp.StatusCode; {
			case sc >= 200 && sc < 400:
				return resp, nil
			case sc >= 400 && sc < 500:
				return nil, fmt.Errorf("failed after %d retries: %s", retry, resp.Status)
			default:
				retryErrs = errors.Join(retryErrs, fmt.Errorf("try %d: %s", retry, resp.Status))
			}

		}

		return nil, fmt.Errorf("failed afdter %d retries: %w", tries, retryErrs)
	}
}

func Trace(rt http.RoundTripper) RoundTripFunc {
	return func(r *http.Request) (*http.Response, error) {
		defer logExec("Trace")()

		traceID, err := uuid.Parse(r.Header.Get("X-Trace-ID"))
		if err != nil {
			traceID = uuid.New()
		}

		trace := trace.Trace{TraceID: traceID, RequestID: uuid.New()}
		ctx := ctxutil.WithValue(r.Context(), trace)

		r = r.WithContext(ctx)

		r.Header.Set("X-Trace-ID", trace.TraceID.String())
		r.Header.Set("X-Request-ID", trace.RequestID.String())

		return rt.RoundTrip(r)
	}
}

func Log(rt http.RoundTripper) RoundTripFunc {
	return func(r *http.Request) (*http.Response, error) {
		defer logExec("Log")()
		var prefix string
		trc, ok := ctxutil.Value[trace.Trace](r.Context())
		if ok {
			prefix = fmt.Sprintf("%s %s: [%s %s]: ", r.Method, r.URL, trc.TraceID, trc.RequestID)
		} else {
			prefix = fmt.Sprintf("%s %s ", r.Method, r.URL)
		}

		logger := log.New(os.Stderr, prefix, log.LstdFlags|log.Lshortfile)
		ctx := ctxutil.WithValue(r.Context(), logger)

		r = r.WithContext(ctx)

		start := time.Now()

		resp, err := rt.RoundTrip(r)
		if err != nil {
			logger.Printf("errored after %s: %s", time.Since(start), err)
			return nil, err
		}

		logger.Printf("%d %s in %s", resp.StatusCode, http.StatusText(resp.StatusCode), time.Since(start))

		return resp, nil
	}
}

func clientMiddleware() http.RoundTripper {
	var rt RoundTripFunc

	const wait, tries = 10 * time.Millisecond, 3

	rt = RetryOn5xx(http.DefaultTransport, wait, tries)
	rt = Log(rt)
	rt = Trace(rt)
	return rt
}
func main() {
	if len(os.Args) < 2 {
		log.Fatal("target url required")
	}

	target := os.Args[1]
	client := &http.Client{Transport: clientMiddleware(), Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(context.TODO(), "GET", target, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
