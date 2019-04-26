package webHandler

import (
	"net/http"
)

type DocServer struct {
	settings WServerSettings
	activeSessions map[string]int
}

type DocumentJob struct {
	uploadPath string
	sessionToken string
	saveDate int64
	jobState string
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

func NewDocServer(srvSettings WServerSettings) DocServer {
	docServer := DocServer{
		settings: srvSettings,
		activeSessions: make(map[string]int),
	}

	return docServer
}

func (dServer *DocServer) Run() error {
	fileServer := http.FileServer(http.Dir("public"))

	http.Handle("/", fileServer)
	http.HandleFunc("/upload", dServer.uploadHandler)
	http.HandleFunc("/job/*", dServer.jobStatus)

	err := http.ListenAndServe(dServer.settings.Ip + ":" + dServer.settings.Port, nil)
	if err != nil {
		return err
	}

	return nil
}

func (dServer *DocServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Store file upload and return a job. (Job: session Token and Document ID)
	/* {
	 *    docId: 0
	 *    sessionKey: asdf1234
	 * }
	 */
}

func (dServer *DocServer) jobStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Return the job status
	/* {
	 *    docId: 0
	 *    jobstate: (processing|success|failed)
	 *    result: "<p>lotsofhtml</p>"
	 * }
	 */
}

