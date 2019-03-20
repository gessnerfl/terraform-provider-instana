package testutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//EchoHandlerFunc is a handler function for the TestHTTPServer which echos the request
//with a http status code 200
func EchoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Write(bytes.NewBufferString("Failed to get request"))
	} else {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write(requestBody)
	}
}

//healthPath is the path for the health endpoint for checking server alive
const healthPath = "/health"

//healthFunc is a handler function which is registered on path /health to check if server
//is running
func healthFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	r.Write(bytes.NewBufferString("OK"))
}

//NewTestHTTPServer create and starts a new TestHTTPServer on random port
func NewTestHTTPServer() *TestHTTPServer {
	router := mux.NewRouter()
	router.HandleFunc(healthPath, healthFunc)
	return &TestHTTPServer{
		router: router,
		port:   RandomPort(),
	}
}

//MinPortNumber the minimum port number used by the http test server
const MinPortNumber = 50000

//MaxPortNumber the maximum port number used by the http test server
const MaxPortNumber = 59000

//RandomPort creates a random port between 50000 and 59000
func RandomPort() int {
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	return random.Intn(MaxPortNumber-MinPortNumber) + MinPortNumber
}

//TestHTTPServer simple helper to mock an http server for testing.
type TestHTTPServer struct {
	router     *mux.Router
	port       int
	httpServer *http.Server
}

//GetPort returns the dynamic server port
func (server *TestHTTPServer) GetPort() int {
	return server.port
}

//AddRoute adds a new route. Routes can only be added before the server was started
func (server *TestHTTPServer) AddRoute(method string, path string, handlerFunc http.HandlerFunc) {
	server.router.HandleFunc(path, handlerFunc).Methods(method)
}

//Start starts the http service with the configured routes
func (server *TestHTTPServer) Start() {
	binding := fmt.Sprintf(":%d", server.port)
	srv := &http.Server{
		Addr:    binding,
		Handler: server.router,
	}
	go func() {
		rootFolder, err := GetRootFolder()
		if err != nil {
			log.Fatalf("Failed to get root folder of project: %s", err)
			return
		}
		certFile := fmt.Sprintf("%s/test-utils/test-server.pem", rootFolder)
		keyFile := fmt.Sprintf("%s/test-utils/test-server.key", rootFolder)
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServeTLS(): %s", err)
		}

	}()
	server.httpServer = srv

	server.waitForServerAlive()
}

func (server *TestHTTPServer) waitForServerAlive() {
	url := fmt.Sprintf("https://localhost:%d/health", server.GetPort())

	for i := 0; i < 5; i++ {
		if resp, err := http.Get(url); err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

//Close stops the http listener
func (server *TestHTTPServer) Close() {
	if server.httpServer != nil {
		server.httpServer.Close()
	}
}
