package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// This was written using the example in the
// "Interfaces and methods" section, as well as
// with the help from this article:
// https://blog.logrocket.com/creating-a-web-server-with-golang/.
type Counter int

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	*ctr++
	fmt.Fprintf(w, "counter = %d\n", *ctr)
}

type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch <- req
	fmt.Fprint(w, "notification sent")
}

func ArgServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args)
}

func myfunc(ch Chan) {
	fmt.Println(<-ch)
}

func main() {
	ctr := new(Counter)
	ch := make(Chan)

	// Create a "goroutine" with the `go` keyword.
	// Had help from this article:
	// https://www.geeksforgeeks.org/goroutines-concurrency-in-golang/.
	// If the main program is terminatd, then so does this one.
	go myfunc(ch)

	// Set up the handlers first before we
	// start listening to a port.
	http.Handle("/", ctr)
	http.Handle("/counter", ctr)
	http.Handle("/args", http.HandlerFunc(ArgServer))
	http.Handle("/test", ch)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
