package main

import (
	"fmt"
	"net/http"

	"github.com/daaku/livereload"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<title>LiveReload</title>
		<script>` + livereload.JS + `</script>
	</head>
	<body>
		<h1>Hello, LiveReload!</h1>
	</body>
</html>
`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/livereload", livereload.Handler)
	fmt.Println("Starting server on http://127.0.0.1:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
