package webserver 


import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"io"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
	"../blockchain"
)

type WebServer struct {
	Blockchain 			*blockchain.Blockchain
}

func NewWebServer() *WebServer {

	newChain := blockchain.NewBlockchain()
	return &WebServer {
		Blockchain:	newChain,
	}
}


func (webserver *WebServer) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", webserver.handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", webserver.handleWriteBlock).Methods("POST")

	return muxRouter
}

func (webserver *WebServer) handleGetBlockchain(w http.ResponseWriter, r *http.Request){
	bytes, err := json.MarshalIndent(webserver.Blockchain, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

type Message struct {
	BPM int
}

func (webserver *WebServer) handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		webserver.responseWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	newBlock, err := webserver.Blockchain.GenerateBlock( m.BPM)
	if err != nil {
		webserver.responseWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}

	if webserver.Blockchain.IsBlockValid(newBlock, webserver.Blockchain.Chain[len(webserver.Blockchain.Chain) - 1]) {
		newBlockchain := append(webserver.Blockchain.Chain, newBlock)
		webserver.Blockchain.ReplaceChain(newBlockchain)
		spew.Dump(webserver.Blockchain.Chain)
	}
	webserver.responseWithJSON(w, r, http.StatusCreated, newBlock)
}

func (webserver *WebServer) responseWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internall Server Error"))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func (webserver *WebServer) Run(httpAddr string) error {
	mux := webserver.makeMuxRouter()
	log.Printf("Listening on %v", httpAddr)

	s := &http.Server{
		Addr:		":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}