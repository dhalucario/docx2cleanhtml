class JobDomController {

    constructor(config) {
        this.uploader = null;
        this.jobList = null;

        if (config.hasOwnProperty('uploader') && typeof config.uploader == 'string') {
            let uploaderElem = document.getElementById(config.uploader);
            if (uploaderElem && uploaderElem instanceof HTMLElement) {
                this.uploader = uploaderElem;
            }
        }

        if (config.hasOwnProperty('jobList') && typeof config.jobList == 'string') {
            let jobListElem = document.getElementById(config.jobList);
            if (jobListElem && jobListElem instanceof HTMLElement) {
                this.jobList = jobListElem;
            }
        }

    }

    initDragDrop(onFiles, onError) {
        if (this.uploader) {
            this.uploader.addEventListener('dragover', function (e) {
                e.stopPropagation();
                e.preventDefault();
                e.dataTransfer.dropEffect = 'copy';
            });

            this.uploader.addEventListener('drop', (e) => {
                e.preventDefault();
                e.stopPropagation();

                let dropFiles = e.dataTransfer.files;

                if (dropFiles) {
                    if (e.dataTransfer.types.includes('Files')) {
                        onFiles(dropFiles)
                    } else {
                        onError('Please only drop files in here');
                    }

                } else {
                    onError('Nothing has been dropped in here');
                }
            });
        }
    }

    addJob(filename) {
        if (this.jobList) {
            let statusList = document.getElementById('status-wrapper');

            if (statusList) {
                let htmlTemplate =
                    '<div class="container">'+
                    '    <div class="row">'+
                    '        <div class="col-12 col-md-4 pr-md-0">'+
                    '            <h3>' + filename + '</h3>'+
                    '        </div>'+
                    '    </div>'+
                    '    <div class="row">'+
                    '        <div class="col-12 col-md-10">'+
                    '            <div class="progress-wrapper">'+
                    '                <progress class="file-progress" value="0" max="100"></progress>'+
                    '            </div>'+
                    '        </div>'+
                    '        <div class="col-12 col-md-2">'+
                    '            <p class="progress-status mb-0">Uploading...</p>'+
                    '        </div>'+
                    '    </div>'+
                    '    <div class="row">'+
                    '        <div class="col-12 textarea-wrapper">'+
                    '            <textarea class="output-area"></textarea>'+
                    '        </div>'+
                    '    </div>';



                let documentElement = document.createElement('div');
                documentElement.dataset.docId = "-1";
                documentElement.innerHTML = htmlTemplate;
                statusList.appendChild(documentElement);

                return documentElement;
            } else {
                return null;
            }
        } else {
            console.warn("Cannot add job, no joblist element available");
            return null;
        }
    }

    showResult(job, result) {
        let elem;

        if (job instanceof HTMLElement) {
            elem = job
        } else if (typeof job === 'string') {
            if (this.jobList) {
                let jobQuery = this.jobList.querySelectorAll('div[data-doc-id="' + job + '"]');
                if (jobQuery && jobQuery.length === 1) {
                    elem = jobQuery[0];
                }
            }
        }

        if (elem) {
            let ouputArea = elem.querySelectorAll('.output-area');
            if (ouputArea && ouputArea.length > 0) {
                ouputArea[0].value = result;
                ouputArea[0].classList.add('show');
            }
        }
    }

    setJobStatus(job, statusText, progress) {
        let elem;

        if (job instanceof HTMLElement) {
            elem = job
        } else if (typeof job === 'string') {
            if (this.jobList) {
                let jobQuery = this.jobList.querySelectorAll('div[data-doc-id="' + job + '"]');
                if (jobQuery && jobQuery.length === 1) {
                    elem = jobQuery[0];
                }
            }
        }
        if (elem) {
            let progressBar = elem.querySelectorAll('progress');
            if (progressBar && progressBar.length > 0) {
                progressBar[0].value = progress;
            }

            let statusElem = elem.querySelectorAll('.progress-status');
            if (statusElem && statusElem.length > 0) {
                statusElem[0].innerHTML = statusText;
            }
        }
    }
}

module.exports = JobDomController;