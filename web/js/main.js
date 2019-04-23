document.addEventListener('DOMContentLoaded', () => {

    let uploaderElement = document.getElementById('uploader');

    if (uploaderElement) {
        uploaderElement.addEventListener('dragover', function (e) {
            e.stopPropagation();
            e.preventDefault();
            e.dataTransfer.dropEffect = 'copy';
        });

        uploaderElement.addEventListener('drop', (dropEvent) => {
            dropEvent.preventDefault();
            dropEvent.stopPropagation();

            let dropFiles = dropEvent.dataTransfer.files;

            if (dropFiles) {
                debugger;
                for (let i = 0; i < dropFiles.length; i++) {
                    if (dropEvent.dataTransfer.types.includes("Files") && dropEvent.dataTransfer.types.length === 1) {
                        checkAndUploadFile(dropFiles[i]).then((session) => {
                            pollJob(session);
                        }).catch((err) => {
                            console.log(err);
                            switch (err) {
                                case "Object is not a file":
                                    alert("Dropped object is not a file");
                                    break;
                                default:
                                    alert("Something went wrong");
                                    break;
                            }
                        });
                    } else {
                        alert("Please only drag files into the area.")
                    }
                }
            }
        });
    } else {
        console.log('Couldn\'t find "uploader"')
    }

});

function checkAndUploadFile(file) {
    return new Promise((resolve, reject) => {
        let session = {
            key: '',
            docId: ''
        };

        sendSimpleRequest('/upload').then((res) => {
            if (res.error) {
                reject(res.error)
            } else {
                session.key = res.key;
                session.docId = res.docId;
                createJobStatus(session.docId, file.name);
                resolve(session);
            }
        }).catch((err) => {
            reject(err);
        });
    });
}

function createJobStatus(docId, filename) {
    let statusList = document.getElementById('status-wrapper');

    if (statusList) {
        let htmlTemplate =
            '<div class="row">' +
            '<div class="col-12 col-md-4 pr-md-0">' +
            '<h3>' + filename + '</h3>' +
            '</div>' +
            '<div class="col-12 col-md-8 pl-md-0">' +
            '<div class="row">' +
            '<div class="col-10 progress-wrapper">' +
            '<progress class="file-progress" value="0" max="100"></progress>' +
            '</div>' +
            '<div class="col-2">' +
            '<p class="progress-status">Uploading...</p>' +
            '</div>' +
            '</div>' +
            '</div>' +
            '</div>' +
            '<div class="row">' +
            '<div class="col-12 textarea-wrapper">' +
            '<textarea class="output-area"></textarea>' +
            '</div>' +
            '</div>';


        let documentElement = document.createElement('div');
        documentElement.innerHTML = htmlTemplate;
        statusList.appendChild(documentElement);
    }

}

function sendSimpleRequest(path, body) {
    return new Promise((resolve, reject) => {
        let req = new XMLHttpRequest();

        req.addEventListener('load', (res) => {
            if (req.status === 200) {
                if (req.responseType === 'json') {
                    try {
                        resolve(JSON.parse(req.responseText))
                    } catch (err) {
                        reject(err)
                    }
                } else {
                    resolve(req.responseText);
                }
            } else {
                reject('Status: ' + req.status)
            }
        });

        req.addEventListener('error', (err) => {
            reject(err);
        });

        // Testing...
        resolve({key: "asdf1234", docId: 1});
        // req.open('POST', window.location.hostname + path, true);
        // req.send(body);
    });
}

function pollJob(session) {

    

    return setTimeout(pollJob, 1000);
}