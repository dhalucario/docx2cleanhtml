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

                jobNetworkController.checkAndUploadFile(file).then((session) => {

                    debugger;

                    jobNetworkController.pollJob(session).then((finishedRes)=>{
                        jobDomController.setJobStatus(statusElement,"Working...", 60)
                    }).catch((err)=>{
                        console.log(err);
                        alert(err);
                    })
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

