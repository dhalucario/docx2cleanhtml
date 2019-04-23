package webHandler

import (
	"log"
	"net/http"
)

type DocServer struct {
	ip   string
	port string
}

func New(ip string, port string) DocServer {
	docServer := DocServer{
		ip,
		port,
	}
	return docServer

}

func (docServer *DocServer) Run() {
	http.HandleFunc("/", docServer.defaultHandler)

	httpErr := http.ListenAndServe(docServer.ip+":"+docServer.port, nil)
	if httpErr != nil {
		log.Fatal(httpErr)
	}
}

func (docServer *DocServer) defaultHandler(w http.ResponseWriter, r *http.Request) {

}
