package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-zoo/bone"
	"github.com/justinas/alice"
)

type Configuration struct {
	ServerHost  string
	ServicePort int
	VCSType     string
	RepoBaseURL string
	DebugOutput bool
	CertFile    string
	KeyFile     string
	Modules     []string
}

type GoPrivateRepoMetaEnpointServer struct {
	config Configuration
	mux    *bone.Mux
}

func (n *GoPrivateRepoMetaEnpointServer) GetConfig() *Configuration {
	return &n.config
}

func (n *GoPrivateRepoMetaEnpointServer) GetMux() *bone.Mux {
	return n.mux
}

// MakeServer creates a fully operational GoPrivateRepoMetaEnpointServer
func MakeServer() *GoPrivateRepoMetaEnpointServer {
	retval := &GoPrivateRepoMetaEnpointServer{mux: bone.New(Serve)}
	retval.InitConfig()
	// init mux routes
	retval.InitMux()
	return retval
}

func (n *GoPrivateRepoMetaEnpointServer) InitConfig() {
	log.Println("Initializing server...")
	// load config file from json
	filename := "config.json"

	// load the configuration from json file
	byts, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Cant open config file: " + filename + "\nerror: " + err.Error())
		return
	}
	err = json.Unmarshal([]byte(byts), &n.config)
	if err != nil {
		log.Fatal("Cant parse config file " + filename + "\nerror: " + err.Error())
		return
	}
	log.Println("Loaded configuration from config.json : \n" + string(byts))
}

func Serve(mux *bone.Mux) *bone.Mux {
	mux.Serve = func(rw http.ResponseWriter, req *http.Request) {
		tr := time.Now()
		mux.DefaultServe(rw, req)
		log.Println("MUX:Serve Req:", req.RemoteAddr, "in", time.Since(tr))
	}
	return mux
}

func (n *GoPrivateRepoMetaEnpointServer) InitMux() {
	handlers := MakeHandlers(n.config.DebugOutput)
	// bundle up default middleware
	_ = alice.New(GetRateLimiterHandler(), handlers.TimeoutHandler,
		handlers.LoggingHandler, handlers.RecoverHandler)
	secureMiddleware := alice.New(GetRateLimiterHandler(), handlers.TimeoutHandler,
		handlers.LoggingHandler, handlers.RecoverHandler)
	n.mux.NotFoundFunc(Handle404)
	// domain handlers
	n.mux.Handle("/:id", secureMiddleware.ThenFunc(http.HandlerFunc(n.GoPrivateRepoMetaEndpointHandler)))
}

func (n *GoPrivateRepoMetaEnpointServer) DoServe() {
	// start mux listening
	servicePort := n.config.ServicePort
	certFile := n.config.CertFile
	keyFile := n.config.KeyFile
	log.Println("GoPrivateRepoMetaEnpointServer service starting at port " + strconv.Itoa(servicePort) + "...")
	if certFile != "" {
		log.Println("starting in secure mode (https)")
		log.Println("using cert file " + certFile)
		log.Println("using key file " + keyFile)
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(servicePort), certFile, keyFile, n.mux))
	} else {
		log.Println("starting in insecure mode (http only) - no certfificate configured")
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(servicePort), n.mux))
	}
}
