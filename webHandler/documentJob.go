package webHandler

import (
	"github.com/google/uuid"
	Document "leong/docx2cleanhtml/simpleDocxParser"
	"log"
	"os"
	"syscall"
)

type JobStatus int

const (
	JobEmpty      JobStatus = 0
	JobError      JobStatus = 1
	JobDone       JobStatus = 2
	JobProcessing JobStatus = 3
)

type DocumentJob struct {
	uploadPath   string
	sessionToken string
	saveDate     int64
	jobState     JobStatus
}

func (dj *DocumentJob) Clear() {
	if dj.jobState != JobEmpty {
		err := os.Remove(dj.uploadPath)
		if err != nil {
			if err != syscall.ENOENT {
				log.Panic(err)
			}
		}
	}
	dj.jobState = JobProcessing
}

func (dj *DocumentJob) InitUUID() {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Panic(err)
	}
	dj.sessionToken = newUUID.String()
}

func (dj *DocumentJob) Status() JobStatus {
	return dj.jobState
}

func (dj *DocumentJob) SetPath(filepath string) {
	dj.uploadPath = filepath
}

func (dj *DocumentJob) ProcessFile() (string, error) {
	doc, err := Document.New(dj.uploadPath)
	if err != nil {
		dj.jobState = JobError
		return "", err
	}

	err = doc.ReadRelations()
	if err != nil {
		dj.jobState = JobError
		return "", err
	}

	return doc.HTML(), nil
}