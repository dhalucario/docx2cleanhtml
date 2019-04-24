package webHandler

import (
	"net/http"
)

type docServer struct {
	settings WServerSettings
	activeSessions map[string]int

}

type WServerSettings struct {
	Ip string
	Port string
}

func (wsrvSettings *WServerSettings)AutocompleteEmpty() {
	if wsrvSettings.Port == "" {
		wsrvSettings.Port = "8080"
	}
	if wsrvSettings.Ip == "" {
		wsrvSettings.Ip = "0.0.0.0"
	}
}

func NewDocServer(srvSettings WServerSettings) docServer {
	docServer := docServer{
		settings: srvSettings,
		activeSessions: make(map[string]int),
	}

	return docServer
}

func (dServer *docServer) Run() error {
	http.HandleFunc("/", dServer.showIndexHandler)
	http.HandleFunc("/upload", dServer.uploadHandler)

	err := http.ListenAndServe(dServer.settings.Ip + ":" + dServer.settings.Port, nil)
	if err != nil {
		return err
	}

	return nil
}

func (dServer *docServer) showIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Just send index.html
}

func (dServer *docServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get file and return a job. (Session Token and Document ID)
}
