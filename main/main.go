package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// hr := hostrouter.New()

	// // Requests to api.domain.com
	// hr.Map("", newHellStudios()) // default
	// hr.Map("newhellstudios.com", newHellStudios())

	// // Requests to doma.in
	// hr.Map("rozeblud.com", rozeBlud())

	// // // Requests to *.doma.in
	// // hr.Map("*.doma.in", rozeBlud())

	// // // Requests to host that isn't defined above
	// // hr.Map("*", everythingElseRouter())

	// // Mount the host router
	// r.Mount("/", hr)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "newhellstudios.com"))
	FileServer(r, "/", filesDir)

	http.ListenAndServe(":80", r)
}

// Router for the API service
func newHellStudios() chi.Router {
	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "newhellstudios.com"))
	FileServer(r, "/", filesDir)
	return r
}

// Router for the Short URL service
func rozeBlud() chi.Router {
	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "rozeblud.com"))
	FileServer(r, "/", filesDir)
	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
