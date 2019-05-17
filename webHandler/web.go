package webHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

type DocServer struct {
	settings          WServerSettings
	activeSessions    map[string]int
	jobController     JobController
	uploadFilesPath   string
	downloadFilesPath string
}

type WServerSettings struct {
	Ip   string
	Port string
}

type StatusRequest struct {
	JobId      int    `json:"jobId"`
	JobSession string `json:"sessionKey"`
}

func (wsrvSettings *WServerSettings) AutocompleteEmpty() {
	if wsrvSettings.Port == "" {
		wsrvSettings.Port = "8080"
	}
	if wsrvSettings.Ip == "" {
		wsrvSettings.Ip = "0.0.0.0"
	}
}

func NewDocServer(srvSettings WServerSettings, cacheDocCount int, inputPath string ,outputPath string) (*DocServer, error) {
	var uploadFilesFolder string
	currentDirectory, err := os.Getwd()

	if outputPath == "" {
		if err != nil {
			panic(err)
		}
		uploadFilesFolder = path.Join(currentDirectory, "uploads")
	} else {
		uploadFilesFolder = outputPath
	}

	if inputPath == "" {
		if err != nil {
			panic(err)
		}
		uploadFilesFolder = path.Join(currentDirectory, "public/output/")
	} else {
		uploadFilesFolder = inputPath
	}

	docJobController, err := NewJobController(cacheDocCount)
	if err != nil {
		return nil, err
	}

	docServer := DocServer{
		settings:        srvSettings,
		activeSessions:  make(map[string]int),
		jobController:   *docJobController,
		uploadFilesPath: uploadFilesFolder,
	}

	return &docServer, nil
}

func (dServer *DocServer) Run() error {
	fileServer := http.FileServer(http.Dir("public"))

	http.Handle("/", fileServer)
	http.HandleFunc("/upload", dServer.uploadHandler)
	http.HandleFunc("/job/status", dServer.jobStatus)

	err := http.ListenAndServe(dServer.settings.Ip+":"+dServer.settings.Port, nil)
	if err != nil {
		return err
	}

	return nil
}

func (dServer *DocServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currentJobID, err := dServer.jobController.InitFreeJob()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savePath := path.Join(dServer.uploadFilesPath, dServer.jobController.jobs[currentJobID].sessionToken) + ".docx"

	saveFile, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}

	requestFile, _, err := r.FormFile("doc_file")
	if err != nil {
		jsonErr, webErr := json.Marshal(map[string]string{
			"err": err.Error(),
		})
		if webErr != nil {
			panic(err)
		}

		_, err = w.Write([]byte(jsonErr))
		if err != nil {
			panic(err)
		}
	}

	_, err = io.Copy(saveFile, requestFile)
	if err != nil {
		panic(err)
	}

	currentJob, err := dServer.jobController.AddJobContent(savePath)
	if err != nil {
		panic(err)
	}

	jsonJob, err := json.Marshal(map[string]string{
		"jobId":      strconv.Itoa(currentJobID),
		"sessionKey": currentJob.sessionToken,
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(jsonJob))
	if err != nil {
		panic(err)
	}

	dServer.processFile(currentJobID)

}

func (dServer *DocServer) jobStatus(w http.ResponseWriter, r *http.Request) {

	// TODO: Return the job status
	/* {
	 *    jobstate: (processing|success|failed)
	 *    result: "<p>lotsofhtml</p>"
	 * }
	 */

	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var rq StatusRequest
	err = json.Unmarshal(b, &rq)
	if err != nil {
		panic(err)
	}

	jobStatus, err := dServer.jobController.StatusWithSession(rq.JobId, rq.JobSession)

	statusResponse := map[string]string{
		"jobState": strconv.Itoa(int(jobStatus)),
		"result":   "",
	}

	if jobStatus == JobDone {
		htmlSaveFile, err := os.Open(path.Join(dServer.downloadFilesPath))
		if err != nil {
			panic(err)
		}

		htmlResult, err := ioutil.ReadAll(htmlSaveFile)
		statusResponse["result"] = string(htmlResult)
	}

	jsonStatus, err := json.Marshal(statusResponse)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(jsonStatus))
	if err != nil {
		panic(err)
	}
}

func (dServer *DocServer) processFile(jobId int) {
	savePath := path.Join(dServer.uploadFilesPath, dServer.jobController.jobs[jobId].sessionToken) + ".html"
	saveFile, err := os.Create(savePath)
	if err != nil {
		fmt.Println(err.Error())
		dServer.jobController.jobs[jobId].jobState = JobError
	}

	htmlContent, err := dServer.jobController.jobs[jobId].ProcessFile()

	_, err = saveFile.Write([]byte(htmlContent))
	if err != nil {
		fmt.Println(err.Error())
		dServer.jobController.jobs[jobId].jobState = JobError
	}

	if err == nil {
		dServer.jobController.jobs[jobId].jobState = JobDone
	}

}
