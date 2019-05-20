package webHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime/debug"
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
	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var uploadFilesFolder string
	var downloadFilesPath string

	if outputPath == "" {
		uploadFilesFolder = path.Join(currentDirectory, "uploads")
	} else {
		uploadFilesFolder = outputPath
	}

	if inputPath == "" {
		downloadFilesPath = path.Join(currentDirectory, "public/output/")
	} else {
		downloadFilesPath = inputPath
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
		downloadFilesPath: downloadFilesPath,
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
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	savePath := path.Join(dServer.uploadFilesPath, dServer.jobController.jobs[currentJobID].sessionToken) + ".docx"

	saveFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	requestFile, _, err := r.FormFile("doc_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	_, err = io.Copy(saveFile, requestFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	currentJob, err := dServer.jobController.AddJobContent(currentJobID, savePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	jsonJob, err := json.Marshal(map[string]interface{}{
		"jobId":      currentJobID,
		"sessionKey": currentJob.sessionToken,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(jsonJob))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	dServer.processFile(currentJobID)

}

func (dServer *DocServer) jobStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	var rq StatusRequest
	err = json.Unmarshal(b, &rq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	jobStatus, err := dServer.jobController.StatusWithSession(rq.JobId, rq.JobSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	statusResponse := map[string]string{
		"jobState": strconv.Itoa(int(jobStatus)),
		"result":   "",
	}

	if jobStatus == JobDone {
		htmlSaveFile, err := os.OpenFile(path.Join(dServer.downloadFilesPath, rq.JobSession + ".html"), os.O_RDONLY, 0775)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			debug.PrintStack()
			return
		}

		htmlResult, err := ioutil.ReadAll(htmlSaveFile)
		statusResponse["result"] = string(htmlResult)
	}

	jsonStatus, err := json.Marshal(statusResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(jsonStatus))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		debug.PrintStack()
		return
	}
}

func (dServer *DocServer) processFile(jobId int) {
	savePath := path.Join(dServer.downloadFilesPath, dServer.jobController.jobs[jobId].sessionToken) + ".html"
	saveFile, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		fmt.Println(err.Error())
		dServer.jobController.jobs[jobId].jobState = JobError
	}

	htmlContent, err := dServer.jobController.jobs[jobId].ProcessFile()
	if err != nil {
		fmt.Println(err.Error())
		dServer.jobController.jobs[jobId].jobState = JobError
	}

	_, err = saveFile.Write([]byte(htmlContent))
	if err != nil {
		fmt.Println(err.Error())
		dServer.jobController.jobs[jobId].jobState = JobError
	}

	if err == nil {
		dServer.jobController.jobs[jobId].jobState = JobDone
	}

}
