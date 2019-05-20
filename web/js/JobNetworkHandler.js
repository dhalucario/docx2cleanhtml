class JobNetworkController {

    checkAndUploadFile(file) {
        return new Promise((resolve, reject) => {
            let session = {
                sessionKey: '',
                jobId: ''
            };

            this.uploadFile('/upload', file).then((res) => {
                if (res.error) {
                    reject(res.error)
                } else {
                    if (res.hasOwnProperty('jobId') && res.hasOwnProperty('sessionKey')) {
                        session.sessionKey = res.sessionKey;
                        session.jobId = res.jobId;
                        resolve(session);
                    } else {
                        reject("Invalid response");
                    }
                }
            }).catch((err) => {
                reject(err);
            });
        });
    }

    sendSimpleRequest(path, body) {
        return new Promise((resolve, reject) => {
            let req = new XMLHttpRequest();

            req.addEventListener('load', () => {
                if (req.status >= 200 && req.status < 300) {
                    try {
                        resolve(JSON.parse(req.responseText))
                    } catch (err) {
                        resolve(req.responseText)
                    }
                } else {
                    reject('Status: ' + req.status)
                }
            });

            req.addEventListener('error', (err) => {
                reject(err);
            });

            req.open('POST', path, true);
            req.send(body);
        });
    }

    uploadFile(path, file) {
        return new Promise((resolve, reject) => {
            let req = new XMLHttpRequest();
            let formData = new FormData();
            formData.append('doc_file', file);

            req.addEventListener('load', () => {
                if (req.status >= 200 && req.status < 300) {
                    try {
                        resolve(JSON.parse(req.responseText))
                    } catch (err) {
                        resolve(req.responseText)
                    }
                } else {
                    reject('Status: ' + req.status)
                }
            });

            req.addEventListener('error', (err) => {
                reject(err);
            });

            req.open('POST', path, true);
            req.send(formData);
        });
    }

    pollJob(session) {
        return new Promise((resolve, reject) => {
            this.sendSimpleRequest('/job/status', JSON.stringify(session)).then((res) => {
                if (res.err) {
                    console.log(res.err);
                } else {
                    if (res.hasOwnProperty('jobState')) {
                        switch (res.jobState) {
                            case '2':
                                res.sessionKey = session.sessionKey;
                                resolve(res);
                                break;
                            case '1':
                                setTimeout(() => {
                                    this.pollJob(session).then((finished) => {
                                        resolve(finished)
                                    }).catch((err) => {
                                        reject(err)
                                    })
                                }, 5000);
                                break;
                            case 'failed':
                                reject(res.err);
                                break;
                        }
                    }
                }
            }).catch((err) => {
                console.log(err)
            });
        })
    }

}

module.exports = JobNetworkController;
