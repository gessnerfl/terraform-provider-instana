package testutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const contentTypeHeaderName = "Content-Type"
const healthPath = "/health"

// EchoHandlerFunc is a handler function for the TestHTTPServer which echos the request
// with a http status code 200
func EchoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err2 := r.Write(bytes.NewBufferString("Failed to get request"))
		if err2 != nil {
			log.Fatalf("failed to write error response; %s", err2)
		}
	} else {
		w.Header().Set(contentTypeHeaderName, r.Header.Get(contentTypeHeaderName))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(requestBody)
		if err != nil {
			log.Fatalf("failed to write response body; %s", err)
		}
	}
}

// healthFunc is a handler function which is registered on path /health to check if server
// is running
func healthFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := r.Write(bytes.NewBufferString("OK"))
	if err != nil {
		log.Fatalf("failed to write health response; %s", err)
	}
}

// NewTestHTTPServer create and starts a new TestHTTPServer on random port
func NewTestHTTPServer() TestHTTPServer {
	router := mux.NewRouter()
	router.HandleFunc(healthPath, healthFunc)
	return &testHTTPServerImpl{
		router:      router,
		callCounter: make(map[string]int),
	}
}

// TestHTTPServer simple helper to mock an http server for testing.
type TestHTTPServer interface {
	GetPort() int
	GetCallCount(method string, path string) int
	AddRoute(method string, path string, handlerFunc http.HandlerFunc)
	Start()
	Close()
	WriteInternalServerError(w http.ResponseWriter, err error)
	WriteJSONResponse(w http.ResponseWriter, jsonData []byte)
}

type testHTTPServerImpl struct {
	router      *mux.Router
	port        *int
	httpServer  *http.Server
	listener    net.Listener
	callCounter map[string]int
}

// GetPort returns the dynamic server port
func (server *testHTTPServerImpl) GetPort() int {
	if server.port == nil {
		return -1
	}
	return *server.port
}

// GetCallCount returns the call counter for the given method and path
func (server *testHTTPServerImpl) GetCallCount(method string, path string) int {
	key := method + "_" + path
	val, ok := server.callCounter[key]
	if !ok {
		return 0
	}
	return val
}

// AddRoute adds a new route. Routes can only be added before the server was started
func (server *testHTTPServerImpl) AddRoute(method string, path string, handlerFunc http.HandlerFunc) {
	server.router.HandleFunc(path, server.wrapHandlerFunc(handlerFunc)).Methods(method)
}

func (server *testHTTPServerImpl) wrapHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + "_" + r.URL.Path
		val, ok := server.callCounter[key]
		if !ok {
			val = 0
		}
		server.callCounter[key] = val + 1
		handlerFunc(w, r)
	}
}

// Start starts the http service with the configured routes
func (server *testHTTPServerImpl) Start() {
	srv := &http.Server{
		Handler:           server.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	binding := fmt.Sprintf(":%d", port)

	go func() {
		rootFolder, err := GetRootFolder()
		if err != nil {
			log.Fatalf("Failed to get root folder of project: %s", err)
			return
		}
		certFile := fmt.Sprintf("%s/testutils/test-server.pem", rootFolder)
		keyFile := fmt.Sprintf("%s/testutils/test-server.key", rootFolder)
		if err = srv.ServeTLS(l, certFile, keyFile); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start http server using binding %s: %s", binding, err)
		}

	}()
	server.httpServer = srv
	server.listener = l
	server.port = &port

	server.waitForServerAlive()
}

func (server *testHTTPServerImpl) waitForServerAlive() {
	url := fmt.Sprintf("https://localhost:%d/health", server.GetPort())

	for i := 0; i < 5; i++ {
		if resp, err := http.Get(url); err == nil && resp.StatusCode == 200 { //nolint:gosec
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Close stops the http listener
func (server *testHTTPServerImpl) Close() {
	if server.httpServer != nil {
		err := server.httpServer.Close()
		if err != nil {
			log.Fatalf("failed to close http server; %s", err)
		}
	}
}

// WriteInternalServerError Writes the provided error message as a response message and sets status code 501 - Internal Server Error with content type text/plain
func (server *testHTTPServerImpl) WriteInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set(contentTypeHeaderName, "text/plain; charset=utf-8")
	_, err2 := w.Write([]byte(err.Error()))
	if err2 != nil {
		log.Fatalf("failed to write internal server error; %s", err2)
	}
}

// WriteJSONResponse Writes the provided data with content type application/json and status code 200 OK to the ResponseWriter
func (server *testHTTPServerImpl) WriteJSONResponse(w http.ResponseWriter, jsonData []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentTypeHeaderName, "application/json; charset=utf-8")
	_, err := w.Write(jsonData)
	if err != nil {
		log.Fatalf("failed to json repsonse; %s", err)
	}
}
