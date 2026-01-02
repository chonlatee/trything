package main

import (
	"chonlatee/myroute/ctxutil"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type route struct {
	pattern *regexp.Regexp
	names   []string

	raw     string
	method  string
	handler http.Handler
}

type Router struct{ routes []route }

func buildRoute(pattern string) (re *regexp.Regexp, names []string, err error) {
	if pattern == "" || pattern[0] != '/' {
		return nil, nil, fmt.Errorf("invalid pattern %s: must begin with '/", pattern)
	}

	var buf strings.Builder

	buf.WriteByte('^')

	for _, f := range strings.Split(pattern, "/")[1:] {
		buf.WriteByte('/')

		if len(f) >= 2 && f[0] == '{' && f[len(f)-1] == '}' {
			trimmed := f[1 : len(f)-1]

			if before, after, ok := strings.Cut(trimmed, ":"); ok {
				names = append(names, before)

				buf.WriteByte('(')
				buf.WriteString(after)
				buf.WriteByte(')')
			} else {
				buf.WriteString(trimmed)
			}
		} else {
			buf.WriteString(regexp.QuoteMeta(f))
		}
	}

	for i := range names {
		for j := i + 1; j < len(names); j++ {
			if names[i] == names[j] {
				return nil, nil, fmt.Errorf("duplicate path parameter %s in %q", names[i], pattern)
			}
		}
	}

	buf.WriteByte('$')
	re, err = regexp.Compile(buf.String())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid regexp %s: %w", buf.String(), err)
	}

	return re, names, nil
}

type PathVars map[string]string

var empty = make(PathVars)

func Vars(ctx context.Context) PathVars {
	v, _ := ctxutil.Value[PathVars](ctx)
	return v
}

func (rt *Router) AddRoute(pattern string, h http.Handler, method string) error {
	re, names, err := buildRoute(pattern)
	if err != nil {
		return err
	}

	rt.routes = append(rt.routes, route{
		raw:     pattern,
		pattern: re,
		names:   names,
		method:  strings.ToUpper(strings.TrimSpace(method)),
		handler: h,
	})

	sort.Slice(rt.routes, func(i, j int) bool {
		return len(rt.routes[i].raw) > len(rt.routes[j].raw) || ((len(rt.routes[i].raw) == len(rt.routes[j].raw)) && rt.routes[i].raw < r.routes[j].raw)
	})

	return nil
}

func pathVars(re *regexp.Regexp, names []string, path string) PathVars {
	matches := re.FindStringSubmatch(path)
	if len(matches) != len(names)+1 {
		panic(fmt.Errorf("programmer error: expected regexp %q to match %q", path, re.String()))
	}

	vars := make(PathVars, len(names))
	for i, match := range matches[1:] {
		vars[names[i]] = match
	}

	return vars
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		if route.pattern.MatchString(r.URL.Path) && (route.method == "" || route.method == r.Method) {
			vars := pathVars(route.pattern, route.names, r.URL.Path)
			ctx := ctxutil.WithValue(r.Context(), vars)
			route.handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}

	http.NotFound(w, r)
}

func main() {

}
