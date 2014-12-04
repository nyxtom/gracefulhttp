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
import "github.com/nyxtom/gracefulhttp"

func restartGraceful() {
	
}

func main() {
	server := gracefulhttp.NewServer(":5000")
	http.HandleFunc("/restart", func (w http.ResponseWriter, req *http.Request) {
	})
	server.ListenAndServe(0)
}
```

# LICENSE

MIT
