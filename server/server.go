package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/atotto/clipboard"
	log "github.com/inconshreveable/log15"
	"github.com/lemonade-command/lemonade/lemon"
	"github.com/skratchdot/open-golang/open"
)

var lineEndingOpt string

//var ra *iprange.Range

func handle_copy(w http.ResponseWriter, req *http.Request) {
	/*remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if !ra.IncludeStr(remoteIP) {
		return
	}*/

	// Read body
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	clipboard.WriteAll(lemon.ConvertLineEnding(string(b), lineEndingOpt))
}

func handle_paste(w http.ResponseWriter, req *http.Request) {
	t, err := clipboard.ReadAll()
	if err == nil {
		io.WriteString(w, t)
	}
}

func handle_open(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	url := q.Get("url")
	open.Run(url)
}

func Serve(c *lemon.CLI, logger log.Logger) error {
	port := c.Port
	//allowIP := c.Allow
	lineEndingOpt = c.LineEnding
	/*ra, err = iprange.New(allowIP)
	if err != nil {
		return err
	}*/

	http.HandleFunc("/copy", handle_copy)
	http.HandleFunc("/paste", handle_paste)
	http.HandleFunc("/open", handle_open)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}
