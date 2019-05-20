package webHandler

import (
	"errors"

	"leong/docx2cleanhtml/mutexHelpers"
)

type JobController struct {
	counter mutexHelpers.MutexCounter
	jobs    []DocumentJob
}

// Job Controller
func NewJobController(maxJobs int) (*JobController, error) {
	if maxJobs < 2 {
		return nil, errors.New("minimum job count has to be greater than one")
	} else {

		jobList := make([]DocumentJob, maxJobs)

		for i := 0; i < len(jobList); i++ {
			jobList[i] = DocumentJob{
				sessionToken: "",
				jobState: JobEmpty,
				uploadPath: "",
				saveDate: 0,
			}
		}

		return &JobController{
			counter: mutexHelpers.MutexCounter{},
			jobs: jobList,
		}, nil
	}
}

func (jc *JobController) AddJobContent(job int, filepath string) (*DocumentJob, error) {
	jc.jobs[job].SetPath(filepath)
	return &(jc.jobs[job]), nil
}

func (jc *JobController) InitFreeJob() (int, error) {
	jc.counter.Lock()
	defer jc.counter.Unlock()

	lastJob := jc.counter.Get()
	freeJob := jc.findUnusedJob(lastJob, len(jc.jobs))

	if freeJob != -1 {
		jc.counter.Set(freeJob)
		jc.jobs[freeJob].Clear()
		jc.jobs[freeJob].InitUUID()
		return freeJob, nil
	}

	freeJob = jc.findUnusedJob(0, lastJob)

	if freeJob != -1 {
		jc.counter.Set(freeJob)
		jc.jobs[freeJob].Clear()
		jc.jobs[freeJob].InitUUID()
		return freeJob, nil
	} else {
		return -1, errors.New("no free jobs")
	}

}

func (jc *JobController) findUnusedJob(offset int, end int) int {
	for i := offset; i < end; i++ {
		if jc.jobs[i].jobState != JobProcessing {
			return i
		}
	}
	return -1
}

func (jc *JobController) StatusWithSession(id int, sessionKey string) (JobStatus, error) {
	if jc.jobs[id].sessionToken == sessionKey {
		return jc.jobs[id].jobState, nil
	} else {
		return JobEmpty, errors.New("session key invalid")
	}
}