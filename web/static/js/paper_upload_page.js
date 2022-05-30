const uploadFormId = "uploadForm"

const paperTopicId = "paperTopic"
const paperDescriptionId = "paperDescription"
const uploaderNameId = "uploaderName"
const uploaderSurnameId = "uploaderSurname"
const paperFileId = "paperFile"

const uploadEndpoint = "/api/v1/paper-upload"
const statusEndpoint = "/api/v1/paper-status"

window.addEventListener("load", function () {
        const form = document.getElementById(uploadFormId)
        form.addEventListener('submit', function (event) {
            console.log("upload form has just been submitted")
            event.preventDefault()
            handlePaperUpload()
        })

        function handlePaperUpload() {
            let http = new XMLHttpRequest();
            http.addEventListener('load', function () {
                console.log('Data sent and response loaded.')
            });

            http.addEventListener('error', function () {
                console.log('Sending data failed.')
            });

            http.onreadystatechange = function () {
                if (http.readyState === XMLHttpRequest.DONE) {
                    console.log('paper upload response status code ' + http.status)
                    if (http.status === 202) {
                        processUploadResults(http.responseText).then(r => console.log(r))
                    } else if (http.status === 200) {
                        displayUploadResults(http.responseText)
                    }
                }
            }

            let formData = new FormData()
            formData.set(paperTopicId, document.getElementById(paperTopicId).value)
            formData.set(paperDescriptionId, document.getElementById(paperDescriptionId).value)
            formData.set(uploaderNameId, document.getElementById(uploaderNameId).value)
            formData.set(uploaderSurnameId, document.getElementById(uploaderSurnameId).value)
            formData.set(paperFileId, document.getElementById(paperFileId).files[0])

            http.open('POST', uploadEndpoint, true)
            http.send(formData)
            console.log('paper upload request was sent')
        }

        async function processUploadResults(paperResultJson) {
            document.getElementById('paperUploadView').style.display = 'none';
            document.getElementById('loadingView').style.display = 'block';

            await new Promise(r => setTimeout(r, 2000));

            let paperResult = JSON.parse(paperResultJson)
            let http = new XMLHttpRequest();
            http.open('GET', statusEndpoint + '?paperId=' + paperResult['id'], true);
            http.onreadystatechange = function () {
                if (http.readyState === XMLHttpRequest.DONE) {
                    console.log('get paper processing status response status code ' + http.status)
                    if (http.status === 200) {
                        displayUploadResults(http.responseText)
                    } else if (http.status === 202) {
                        processUploadResults(http.responseText)
                    }
                }
            }
            http.send(null);
            console.log('get paper processing status request was sent')
        }

        function displayUploadResults(responseBody) {
            document.getElementById('loadingView').style.display = 'none';
            document.getElementById('uploadResult').style.display = 'block';

            let uploadResult = JSON.parse(responseBody)
            console.log('paper ' + uploadResult['id'] + ' processing finished with status ' + uploadResult['status'])
            document.getElementById('resultingNFW').innerHTML = uploadResult['NFT']
        }
    });