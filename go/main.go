package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/zserge/lorca"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine.

//  message for increment counter
type counterAddOp struct {
	val  int // increment value +1 or -1
	resp chan bool
}

// message for refreshing UI with latest counter
type refreshCounterOp struct {
	counter int
	resp    chan bool
}

func main() {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 480, 320, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Creates two channels, one for sending increment, and another to refreshUI with latest counter
	counterAddChannel := make(chan counterAddOp)
	refreshCounterChannel := make(chan refreshCounterOp)

	// counter goroutine
	go func() {
		var count = 0
		for {
			select {
			case inc := <-counterAddChannel:
				count = count + inc.val
				ui.Eval(fmt.Sprintf("console.log('Go: inside goroutine for counter',%d)", count))
				op := refreshCounterOp{
					counter: count,
					resp:    make(chan bool)}
				refreshCounterChannel <- op
				<-op.resp
				inc.resp <- true
			}
		}
	}()

	// refresh UI with lastest goroutine
	go func() {
		for {
			select {
			case refresh := <-refreshCounterChannel:
				ui.Eval(fmt.Sprintf("console.log('Go: inside goroutine for refresh',%d)", refresh.counter))
				jsAction := fmt.Sprintf("updateUI(%d);", refresh.counter)
				ui.Eval(jsAction)
				refresh.resp <- true
			}
		}
	}()

	// Create and bind Go object to the UI
	ui.Bind("counterAdd", func(inc int) bool {
		op := counterAddOp{
			val:  inc,
			resp: make(chan bool)}
		counterAddChannel <- op
		<-op.resp
		ui.Eval(`console.log("Go: sent counterAdd message");`)
		return true
	})

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
