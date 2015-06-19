package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.sakura.ad.jp/h-nakamura/intra-doc-search"

	log "github.com/Sirupsen/logrus"
	"github.com/doloopwhile/logrusltsv"
	"github.com/gocraft/web"
	"github.com/hnakamur/gocraft-web-logrus"
)

type Context struct {
	HelloCount int
}

func (c *Context) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

func (c *Context) SayHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, strings.Repeat("Hello ", c.HelloCount), "World!")
}

var (
	logFileName   string
	listenAddress string
)

func init() {
	flag.StringVar(&logFileName, "logfilename", "-", "log file (- for stdout)")
	flag.StringVar(&listenAddress, "listenaddress", "localhost:3000", "listen address (host:port)")
}

func setupLogging(logFileName string) (logFile *os.File, err error) {
	log.SetFormatter(&logrusltsv.Formatter{})
	if logFileName == "-" {
		logFile = os.Stdout
	} else {
		logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	}
	log.SetOutput(logFile)
	return logFile, nil
}

func main() {
	flag.Parse()

	logFile, err := intra_doc_search.SetupLogging(logFileName)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	router := web.New(Context{}). // Create your router
					Middleware(gocraft_web_logrus.LoggerMiddlewareFactory("processed request")). // Use some included middleware
					Middleware(web.ShowErrorsMiddleware).                                        // ...
					Middleware((*Context).SetHelloCount).                                        // Your own middleware!
					Get("/", (*Context).SayHello)                                                // Add a route
	http.ListenAndServe(listenAddress, router) // Start the server!
}
