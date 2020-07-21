// Program shortcut runs at https://cs.dev and lets you quickly search code on Sourcegraph by typing
// `cs.dev/my query` into your browser's URL bar.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

const queryPlaceholder = "$QUERY"

var (
	httpListenAddr = flag.String("http", ":"+getenvOrDefault("PORT", "3980"), "HTTP listen address")
	tlsCertPath    = flag.String("tls-cert", "", "path to TLS certificate file")
	tlsKeyPath     = flag.String("tls-key", "", "path to TLS key file")
	redirectURL    = flag.String("redirect-url", getenvOrDefault("REDIRECT_URL", "https://sourcegraph.com/search?q="+queryPlaceholder+"&patternType=structural"), "URL of redirect destination (use "+queryPlaceholder+" as placeholder in query param)")
	docsURL        = flag.String("docs-url", getenvOrDefault("DOCS_URL", "https://sourcegraph.com"), "URL for docs (when HTTP request path is empty)")
	accessLog      = flag.Bool("access-log", false, "print HTTP access log to stdout")
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	flag.Parse()

	makeRedirectURL, err := parseURLPattern(*redirectURL)
	if err != nil {
		log.Fatal(err)
	}

	var handler http.Handler = &handler{
		makeRedirectURL: makeRedirectURL,
		docsURL:         *docsURL,
	}
	if *accessLog {
		handler = handlers.CombinedLoggingHandler(os.Stderr, handler)
	}
	log.Printf("# Listening on %s", *httpListenAddr)
	if *tlsCertPath != "" && *tlsKeyPath != "" {
		log.Fatal(http.ListenAndServeTLS(*httpListenAddr, *tlsCertPath, *tlsKeyPath, handler))
	} else {
		log.Fatal(http.ListenAndServe(*httpListenAddr, handler))
	}
}

type handler struct {
	makeRedirectURL func(query string) string
	docsURL         string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		http.Redirect(w, r, h.docsURL, http.StatusFound)
		return
	}

	http.Redirect(w, r, h.makeRedirectURL(strings.TrimPrefix(r.URL.Path, "/")), http.StatusFound)
}

func parseURLPattern(urlStr string) (makeURL func(query string) string, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var interpolateKey string
	for key, val := range u.Query() {
		if len(val) != 1 {
			continue
		}
		if strings.Contains(val[0], queryPlaceholder) {
			interpolateKey = key
			break
		}
	}
	if interpolateKey == "" {
		return nil, fmt.Errorf("URL %q does not contain query parameter with query placeholder %q", urlStr, queryPlaceholder)
	}

	return func(query string) string {
		q := u.Query()
		q.Set(interpolateKey, strings.Replace(q.Get(interpolateKey), queryPlaceholder, query, 1))
		return u.ResolveReference(&url.URL{RawQuery: q.Encode()}).String()
	}, nil
}

func getenvOrDefault(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defaultValue
}
