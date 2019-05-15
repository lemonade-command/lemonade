package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/atotto/clipboard"
	log "github.com/inconshreveable/log15"
	"github.com/lemonade-command/lemonade/lemon"
	"github.com/pocke/go-iprange"
	"github.com/skratchdot/open-golang/open"
)

var logger log.Logger
var lineEnding string
var ra *iprange.Range
var port int
var path = "./files"

func handleCopy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Copy only support post", 404)
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	text := lemon.ConvertLineEnding(string(b), lineEnding)
	logger.Debug("Copy:", "text", text)
	clipboard.WriteAll(text)
}

func handlePaste(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Paste only support get", 404)
		return
	}

	t, err := clipboard.ReadAll()
	if err == nil {
		io.WriteString(w, t)
	}
	logger.Debug("Paste: ", "text", t)
}

func translateLoopbackIP(uri string, remoteIP string) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return uri
	}

	host, port, err := net.SplitHostPort(parsed.Host)

	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return uri
	}

	if len(port) == 0 {
		parsed.Host = remoteIP
	} else {
		parsed.Host = fmt.Sprintf("%s:%s", remoteIP, port)
	}

	return parsed.String()
}

func handleOpen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Open only support get", 404)
		return
	}

	q := r.URL.Query()
	uri := q.Get("uri")
	isBase64 := q.Get("base64")
	if isBase64 == "true" {
		decodeURI, err := base64.URLEncoding.DecodeString(uri)
		if err != nil {
			logger.Error("base64 decode error", "uri", uri)
			return
		}
		uri = string(decodeURI)
	}

	transLoopback := q.Get("transLoopback")
	if transLoopback == "true" {
		remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		uri = translateLoopbackIP(uri, remoteIP)
	}

	logger.Info("Open: ", "uri", uri)
	open.Run(uri)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Upload only support post", 404)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Error Retrieving the File", 500)
		logger.Error("Error Retrieving the File", "err", err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error Read the File", 500)
		logger.Error("Error Read the File", "err", err)
		return
	}

	ioutil.WriteFile(path+"/"+handler.Filename, fileBytes, os.ModePerm)

	q := r.URL.Query()
	isOpen := q.Get("open")
	if isOpen == "true" {
		uri := fmt.Sprintf("http://127.0.0.1:%d/files/%s", port, handler.Filename)
		logger.Info("Open: ", "uri", uri)
		open.Run(uri)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			http.Error(w, "Not support method.", 404)
			return
		}

		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "RemoteAddr error.", 500)
			return
		}
		if !ra.IncludeStr(remoteIP) {
			http.Error(w, "Not allow ip.", 503)
			logger.Info("not in allow ip. from: ", remoteIP)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Serve(c *lemon.CLI, _logger log.Logger) error {
	logger = _logger
	lineEnding = c.LineEnding
	port = c.Port

	var err error
	ra, err = iprange.New(c.Allow)
	if err != nil {
		logger.Error("allowIp error")
		return err
	}

	os.MkdirAll(path, os.ModePerm)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(path))))
	http.Handle("/copy", middleware(http.HandlerFunc(handleCopy)))
	http.Handle("/paste", middleware(http.HandlerFunc(handlePaste)))
	http.Handle("/open", middleware(http.HandlerFunc(handleOpen)))
	http.Handle("/upload", middleware(http.HandlerFunc(handleUpload)))
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
	if err != nil {
		return err
	}
	return nil
}
