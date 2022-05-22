window.addEventListener("load", function (){
    const form = document.getElementById( "uploadForm" );
    form.addEventListener( 'submit', function ( event ) {
        event.preventDefault();
        handlePaperUpload();
    } );

    const uploaderName = document.getElementById( "uploaderName" );
    const paperTopic = document.getElementById( "paperTopic" );
    const paperFile = {
        dom    : document.getElementById( "paperFile" ),
        binary : null
    };
    const reviewFile = {
        dom    : document.getElementById( "reviewFile" ),
        binary : null
    };

    const paperFileReader = new FileReader();
    paperFileReader.addEventListener( "load", function () {
        paperFile.binary = paperFileReader.result;
    } );

    if( paperFile.dom.files[0] ) {
        paperFileReader.readAsBinaryString( paperFile.dom.files[0] );
    }

    paperFile.dom.addEventListener( "change", function () {
        if( paperFileReader.readyState === FileReader.LOADING ) {
            paperFileReader.abort();
        }
        paperFileReader.readAsBinaryString( paperFile.dom.files[0] );
    } );

    const reviewFileReader = new FileReader();
    reviewFileReader.addEventListener( "load", function () {
        reviewFile.binary = reviewFileReader.result;
    } );

    if( reviewFile.dom.files[0] ) {
        reviewFileReader.readAsBinaryString( reviewFile.dom.files[0] );
    }

    reviewFile.dom.addEventListener( "change", function () {
        if( reviewFileReader.readyState === FileReader.LOADING ) {
            reviewFileReader.abort();
        }
        reviewFileReader.readAsBinaryString( reviewFile.dom.files[0] );
    } );


    function handlePaperUpload() {
        if( !paperFile.binary && paperFile.dom.files.length > 0 ) {
            setTimeout( handlePaperUpload, 10 );
            return;
        }
        if( !reviewFile.binary && reviewFile.dom.files.length > 0 ) {
            setTimeout( handlePaperUpload, 10 );
            return;
        }

        const boundary = "blob";
        let data = "";

        if ( paperFile.dom.files[0] ) {
            data += "--" + boundary + "\r\n";
            data += 'content-disposition: form-data; '
                + 'name="'         + paperFile.dom.name          + '"; '
                + 'filename="'     + paperFile.dom.files[0].name + '"\r\n';

            data += 'Content-Type: ' + paperFile.dom.files[0].type + '\r\n';
            data += '\r\n';
            data += paperFile.binary + '\r\n';
        }

        if ( reviewFile.dom.files[0] ) {
            data += "--" + boundary + "\r\n";
            data += 'content-disposition: form-data; '
                + 'name="'         + reviewFile.dom.name          + '"; '
                + 'filename="'     + reviewFile.dom.files[0].name + '"\r\n';

            data += 'Content-Type: ' + reviewFile.dom.files[0].type + '\r\n';
            data += '\r\n';
            data += reviewFile.binary + '\r\n';
        }

        data += "--" + boundary + "\r\n";
        data += 'content-disposition: form-data; name="' + uploaderName.name + '"\r\n';
        data += '\r\n';
        data += uploaderName.value + "\r\n";

        data += "--" + boundary + "\r\n";
        data += 'content-disposition: form-data; name="' + paperTopic.name + '"\r\n';
        data += '\r\n';
        data += paperTopic.value + "\r\n";

        data += "--" + boundary + "--";

        let http = new XMLHttpRequest();
        http.onreadystatechange = function () {
            if ( http.readyState === XMLHttpRequest.DONE ) {
                if (http.status === 202) {
                    processUploadResults(JSON.parse(http.responseText)).then(r => console.log(r))
                }
            }
        }

        http.addEventListener( 'load', function() {
            alert( 'Data sent and response loaded.' );
        } );

        http.addEventListener( 'error', function() {
            alert( 'Sending data failed.' );
        } );

        http.open( 'POST', '/paper-upload', true )
        http.setRequestHeader( 'Content-Type', 'multipart/form-data; boundary=' + boundary );
        http.send(data)
    }
    
    async function processUploadResults(sessionInProgressId) {
        document.getElementById('paperUploadView').style.display = 'none';
        document.getElementById('loadingView').style.display = 'block';

        await new Promise(r => setTimeout(r, 2000));

        let http = new XMLHttpRequest();
        http.open('GET', '/paper-upload/status?paperId='+sessionInProgressId, true);
        http.onreadystatechange = function() {
            if ( http.readyState === XMLHttpRequest.DONE ) {
                if ( http.status === 200 ) {
                    displayUploadResults(http.responseText)
                } else if ( http.status === 204 ) {
                    processUploadResults(sessionInProgressId)
                }
            }
        }
        http.send(null);
    }

    function displayUploadResults(responseBody) {
        document.getElementById('loadingView').style.display = 'none';
        document.getElementById('uploadResult').style.display = 'block';

        let uploadResult = JSON.parse(responseBody)
        document.getElementById('resultingNFW').innerHTML = uploadResult['NFT']
    }
});