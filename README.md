# gracefulhttp

gracefulhttp is a simple http server that wraps the http server to handle
graceful restarts as referenced and implemented on
http://grisha.org/blog/2014/06/03/graceful-restart-in-golang/.

View the [docs](http://godoc.org/github.com/nyxtom/gracefulhttp).

## Installation

```
$ go get github.com/nyxtom/gracefulhttp
```

## Example

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/nyxtom/gracefulhttp"
)

func restart(server *gracefulhttp.Server) error {
	fd := server.Fd()
	cmd := exec.Command(os.Args[0], []string{"-fd", fmt.Sprintf("%d", fd)}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func handleFunc(prefix string, fns ...func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(prefix, func(w http.ResponseWriter, req *http.Request) {
		for _, fn := range fns {
			fn(w, req)
		}
	})
}

func logReq(w http.ResponseWriter, req *http.Request) {
	log.Printf("%v %v from %v", req.Method, req.URL, req.RemoteAddr)
}

func main() {
	var fd = flag.Int("fd", 0, "")
	flag.Parse()

	server := gracefulhttp.NewServer(":5000", *fd)
	handleFunc("/", logReq, func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello world\n")
	})
	handleFunc("/restart", logReq, func(w http.ResponseWriter, req *http.Request) {
		err := restart(server)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, req, "/", http.StatusFound)
	})
	handleFunc("/shutdown", logReq, func(w http.ResponseWriter, req *http.Request) {
		server.Close()
	})

	pid := os.Getpid()
	log.SetPrefix(fmt.Sprintf("\033[36m[%d]\033[m ", pid))
	if *fd != 0 {
		log.Printf("listening on existing file descriptor %d\n", *fd)
	} else {
		log.Printf("listening on %s\n", server.Addr)
	}
	server.ListenAndServe()
}
```

# LICENSE

MIT
