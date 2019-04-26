document.addEventListener('DOMContentLoaded', () => {
    let jobDomController = new JobDomController({
        uploader: "uploader",
        jobList: "status-wrapper"
    });
    let jobNetworkConroller = new JobNetworkConroller();

    jobDomController.initDragDrop((dropFiles) => {
        for (let i = 0; i < dropFiles.length; i++) {
            ((file) => {
                let statusElement = jobDomController.addJob(file.name);

                jobNetworkConroller.checkAndUploadFile(file).then((session) => {

                    jobNetworkConroller.pollJob(session).then((finishedRes)=>{
                        jobDomController.setJobStatus(statusElement,"Working...", 60)
                    }).catch((err)=>{
                        console.log(err);
                        alert(err);
                    })
                }).catch((err) => {
                    console.log(err);
                    alert(err);
                });
            })(file)
        }
    }, (err)=>{
        alert("Please only drag files into the drop zone");
        console.log(err);
    });
});

