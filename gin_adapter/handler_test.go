package gin_adapter_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IRelaxxx/servefiles/v3/gin_adapter"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func ExampleHandlerFunc() {
	// This is a webserver using the asset handler provided by
	// github.com/IRelaxxx/servefiles/v3, which has enhanced
	// HTTP expiry, cache control, compression etc.
	// 'Normal' bespoke handlers are included as needed.

	// where the assets are stored (replace as required)
	localPath := "./assets"

	// how long we allow user agents to cache assets
	// (this is in addition to conditional requests, see
	// RFC7234 https://tools.ietf.org/html/rfc7234#section-5.2.2.8)
	maxAge := time.Hour

	h := gin_adapter.NewAssetHandler(localPath).
		WithMaxAge(maxAge).
		WithNotFound(http.NotFoundHandler()). // supply your own
		StripOff(1).
		HandlerFunc("filepath")

	router := gin.Default()
	// ... add other routes / handlers / middleware as required
	router.GET("/files/*filepath", h)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func TestHandlerFunc(t *testing.T) {
	g := NewGomegaWithT(t)

	maxAge := time.Hour
	fs := afero.NewMemMapFs()
	fs.MkdirAll("/foo/bar", 0755)
	afero.WriteFile(fs, "/foo/bar/x.txt", []byte("hello"), 0644)

	h := gin_adapter.NewAssetHandlerFS(fs).
		WithMaxAge(maxAge).
		WithNotFound(http.NotFoundHandler()). // supply your own
		StripOff(1).
		HandlerFunc("filepath")

	router := gin.Default()
	// ... add other routes / handlers / middleware as required
	router.GET("/files/*filepath", h)
	router.HEAD("/files/*filepath", h)

	r, _ := http.NewRequest(http.MethodGet, "http://localhost/files/101/foo/bar/x.txt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	g.Expect(w.Header().Get("Expires")).NotTo(Equal(""))
	g.Expect(w.Body.Len()).To(Equal(5))

	r, _ = http.NewRequest(http.MethodHead, "http://localhost/files/101/foo/bar/x.txt", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	g.Expect(w.Code).To(Equal(200))
	g.Expect(w.Header().Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
	g.Expect(w.Header().Get("Expires")).NotTo(Equal(""))
	g.Expect(w.Body.Len()).To(Equal(0))

	r, _ = http.NewRequest(http.MethodHead, "http://localhost/files/101/foo/baz.png", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	g.Expect(w.Code).To(Equal(404))
}
