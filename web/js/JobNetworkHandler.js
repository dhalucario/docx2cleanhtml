class JobNetworkConroller {


    constructor(baseurl) {
        if (baseurl && baseurl !== "") {
            this.baseurl = baseurl
        } else {
            this.baseurl = window.location.protocol + "://" + window.location.hostname
        }
    }

    checkAndUploadFile(file) {
        return new Promise((resolve, reject) => {
            let session = {
                id: '',
                docId: ''
            };

            this.sendSimpleRequest('/upload').then((res) => {
                if (res.error) {
                    reject(res.error)
                } else {
                    session.id = res.id;
                    session.docId = res.docId;
                    resolve(session);
                }
            }).catch((err) => {
                reject(err);
            });
        });
    }

    sendSimpleRequest(path, body) {
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

            // TODO: Get key and docid from server.
            resolve({key: "asdf1234", docId: 1});
            // req.open('POST', this.baseurl + path, true);
            // req.send(body);
        });
    }

    pollJob(session) {
        return new Promise((resolve, reject)=>{
            this.sendSimpleRequest('/job/' + session.docId).then((res) => {
                if (res.err) {
                    console.log(res.err);
                } else {
                    if (res.hasOwnProperty('jobstate') && typeof res.jobstate == 'string') {
                        switch (res.jobstate) {
                            case 'success':
                                    resolve(res);
                                break;
                            case 'processing':
                                setTimeout(() => {
                                    this.pollJob(session).then((finished)=>{
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