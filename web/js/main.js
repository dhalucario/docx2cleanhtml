const JobNetworkController = require('./JobNetworkHandler');
const JobDomController = require('./JobDomController');

document.addEventListener('DOMContentLoaded', () => {

    let jobDomController = new JobDomController({
        uploader: "uploader",
        jobList: "status-wrapper"
    });
    let jobNetworkController = new JobNetworkController();
    jobDomController.initDragDrop((dropFiles) => {

        for (let i = 0; i < dropFiles.length; i++) {

            ((file) => {
                let statusElement = jobDomController.addJob(file.name);
                jobDomController.setJobStatus(statusElement, 'Uploading...', 33);

                jobNetworkController.checkAndUploadFile(file).then((session) => {
                    jobDomController.setJobStatus(statusElement, 'Processing...', 66);

                    jobNetworkController.pollJob(session).then((finishedRes)=>{
                        jobDomController.setJobStatus(statusElement,'<a href="/output/' + finishedRes.sessionKey + '.html">Done!</a>', 100);
                        jobDomController.showResult(statusElement, finishedRes.result);

                    }).catch((err)=>{
                        console.log(err);
                        alert(err);
                    });

                }).catch((err) => {
                    console.log(err);
                    alert(err);
                });
            })(dropFiles[i])

        }
    }, (err)=>{
        alert("Please only drag files into the drop zone");
        console.log(err);
    });
});

