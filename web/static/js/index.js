const displayStyleHide = 'none'
const displayStyleShow = 'block'

let authorsCnt = 1
const authorNameBaseId = 'authorName'
const authorSurnameBaseId = 'authorSurname'
const authorDegreeBaseId = 'authorDegree'

// available endpoints to cover functionality
const addPaperEndpoint       = '/api/v1/addPaper'
const getPaperStatusEndpoint = '/api/v1/getPaperStatus'
const getPaperInfoEndpoint   = '/api/v1/getPaperInfo'
const searchForPaperEndpoint = '/api/v1/searchForPaper'

// get requests parameters
const paperIdKey       = "paperId"
const paperNftKey      = "paperNft"
const searchPayloadKey = "searchPayload"

// possible response states
const okayState          = "okay"
const processingState    = "processing"
const failedState        = "failed"
const noResultsState     = "noResults"
const lowUniquenessState = "lowUniqueness"

function generateAuthorFields() {
    let nameDiv = document.getElementById('nameDiv')
    let surnameDiv = document.getElementById('surnameDiv')
    let degreeDiv = document.getElementById('degreeDiv')

    let nameLabel = document.getElementById('nameLabel')
    let surnameLabel = document.getElementById('surnameLabel')
    let degreeLabel = document.getElementById('degreeLabel')

    nameDiv.innerHTML=""
    surnameDiv.innerHTML=""
    degreeDiv.innerHTML=""

    nameLabel.setAttribute('for', '')
    surnameDiv.setAttribute('for', '')
    degreeDiv.setAttribute('for', '')

    for (let i = 1; i <= authorsCnt; i++) {
        nameDiv.innerHTML+= `<input type="text" placeholder="Enter ${i} name" id="${authorNameBaseId+i}" required">`
        surnameDiv.innerHTML+= `<input type="text" placeholder="Enter ${i} surname" id="${authorSurnameBaseId+i}" required">`
        degreeDiv.innerHTML+= `<input class="degreeInput" type="text" placeholder="Enter ${i} degree" id="${authorDegreeBaseId+i}" required">`

        nameLabel.setAttribute('for', nameLabel.getAttribute('for') + `${authorNameBaseId+i} `)
        surnameDiv.setAttribute('for', surnameLabel.getAttribute('for') + `${authorSurnameBaseId+i} `)
        degreeDiv.setAttribute('for', degreeLabel.getAttribute('for') + `${authorDegreeBaseId+i} `)
    }
}

function addAuthor() {
    if (authorsCnt === 3) {
        return
    }
    if (authorsCnt === 2) {
        document.getElementById('addAuthorButton').style.display = displayStyleHide
    }

    authorsCnt++
    generateAuthorFields()
}

function showAppPaperForm() {
    console.log("Show add paper form")

    document.getElementById('appPaperFormContainer').style.display = displayStyleShow;
    document.getElementById('getPaperFormContainer').style.display = displayStyleHide;
    document.getElementById('searchForPaperFormContainer').style.display = displayStyleHide;

    document.getElementById('paperUploadView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('addPaperResultView').style.display = displayStyleHide;

    generateAuthorFields()
}

function hideAppPaperForm(){
    document.getElementById('appPaperFormContainer').style.display = displayStyleHide;
}
function hideAddResultView(){
    document.getElementById('addResultDiv').style.display = displayStyleHide;
}

function showGetPaperForm() {
    console.log("Show get paper form")

    document.getElementById('appPaperFormContainer').style.display = displayStyleHide;
    document.getElementById('getPaperFormContainer').style.display = displayStyleShow;
    document.getElementById('searchForPaperFormContainer').style.display = displayStyleHide;

    document.getElementById('paperNftEnterView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('getPaperResultView').style.display = displayStyleHide;
}

function hideGetPaperForm(){
    document.getElementById('getPaperFormContainer').style.display = displayStyleHide;

}

function hideGetResultView(){
    document.getElementById('getResultDiv').style.display = displayStyleHide;

}

function showSearchForPaperForm() {
    console.log("Show search for paper form")

    document.getElementById('appPaperFormContainer').style.display = displayStyleHide;
    document.getElementById('getPaperFormContainer').style.display = displayStyleHide;
    document.getElementById('searchForPaperFormContainer').style.display = displayStyleShow;

    document.getElementById('searchTextEnterView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('searchForPaperResultView').style.display = displayStyleHide;
}

function hideSearchPaperForm(){
    document.getElementById('searchForPaperFormContainer').style.display = displayStyleHide;

}
function hideSearchResultView(){
    document.getElementById('searchResultDiv').style.display = displayStyleHide;
}

async function checkPaperStatus(response) {
    hideAppPaperForm()
    document.getElementById('paperUploadView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleShow;
    document.getElementById('addPaperResultView').style.display = displayStyleHide;

    await new Promise(r => setTimeout(r, 2000));

    let http = new XMLHttpRequest();
    http.onreadystatechange = function () {
        addPaperResponseHandler(http)
    }

    http.open('GET', `${getPaperStatusEndpoint}?${paperIdKey}=${response['id']}`, true);
    http.send(null);
    console.log('get paper status request was sent')
}

function handleAddPaperResult(response) {
    document.getElementById('paperUploadView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('addResultDiv').style.display = displayStyleShow;

    console.log(`paper ${response['id']} upload finished`)
    switch (response['state']) {
        case okayState:
            showUploadSuccessText(response)
            document.getElementById('addPaperResultView').style.display = displayStyleShow;
            break
        case lowUniquenessState:
            notUniqueEnoughResponse(response)
            document.getElementById('addPaperResultView').style.display = displayStyleShow;
            break
        default:
            console.log(`State: ${response['state']}. Message:  with ${response['message']}`)
            showError(response, 'addPaperResultView')
    }
}

function showUploadSuccessText(response) {
    document.getElementById('addPaperResultView').innerHTML =
        `<h3>Paper uploaded successfully</h3>` +
        `<label id="resultNftAddress"><b>Your NFT:</b> ${response['nft']['address']}</label>` +
        `<br><label id="resultNftSymbol"><b>Your NFT symbol:</b> ${response['nft']['symbol']}</label>` +
        `<br><label id="resultNftName"><b>Your genearted NFT:</b> ${response['nft']['name']}</label>` +
        `<br><label id="resultPaperUniqueness"><b>Your paper uniqueness:</b> ${response['uniqueness']}</label>` +
        `<br><label id="resultPaperIpfsHash"><b>Your paper ipfs hash:</b> <a href="${response['ipfsHash']}">Your paper at IPFS</a></label>` +
        `<br><label id="resultPaperIpfsHash"><b>Your nft recovery phrase:</b> ${response['nftRecoveryPhrase']}</label>`
}

function notUniqueEnoughResponse(response) {
    let similarNftHtml = ""
    response['similarPapersNft'].forEach(function(nft) {
        similarNftHtml += `<br><a href="/api/v1${nft}">${nft}</a>`
    })
    document.getElementById('addPaperResultView').innerHTML =
        `<h3>Paper uploaded failed</h3>` +
        `<label id="resultMessage"><b>Message:</b> ${response['message']}</label>` +
        `<br><label id="resultPaperUniqueness"><b>Paper uniqueness:</b> <p style="color:#ff0044">${response['uniqueness']}</p></label>` +
        `<label id="similarPapersNft"><b>Similar papers NFT:</b> ${similarNftHtml}</label>`
}

function addPaperResponseHandler(http) {
    if (http.readyState === XMLHttpRequest.DONE) {
        console.log('response status code ' + http.status)
        if (http.status !== 200 && http.status !== 202) {
            console.log("get paper status failed.")
            return
        }

        let response = JSON.parse(http.responseText)
        if (response['state'] === processingState) {
            checkPaperStatus(response).then(r => console.log(r))
        } else {
            handleAddPaperResult(response)
        }
    }
}
let uploadForm = document.getElementById('uploadForm');
uploadForm.addEventListener('submit', function (event) {
    event.preventDefault()

    let http = new XMLHttpRequest();
    http.addEventListener('load', function () {console.log('Data sent and response loaded')
    uploadForm.reset()
    });
    http.addEventListener('error', function () {console.log('Sending data failed')});
    http.onreadystatechange = function () {
        addPaperResponseHandler(http)
    }

    let formData = new FormData()
    formData.set('paperFile', document.getElementById('paperFile').files[0])
    formData.set('paperTopic', document.getElementById('paperTopic').value)
    formData.set('paperDescription', document.getElementById('paperDescription').value)
    for (let i = 1; i <= authorsCnt; i++) {
        formData.set(`${authorNameBaseId+i}`, document.getElementById(`${authorNameBaseId+i}`).value)
        formData.set(`${authorSurnameBaseId+i}`, document.getElementById(`${authorSurnameBaseId+i}`).value)
        formData.set(`${authorDegreeBaseId+i}`, document.getElementById(`${authorDegreeBaseId+i}`).value)
    }

    http.open('POST', addPaperEndpoint, true)
    http.send(formData)

    console.log('add paper request was sent')
})

function handleGetPaperInfoResult(response) {
    document.getElementById('paperNftEnterView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('getResultDiv').style.display = displayStyleShow;

    switch (response['state']) {
        case okayState:
            document.getElementById('getPaperResultView').innerHTML =
                `<h3>Paper data for NFT ${response['nft']}</h3>` +
                `<label id="resultTopic"><b>Topic</b>: ${response['metadata']['topic']}</label>` +
                `<br><label id="resultDescription"><b>Description</b>: ${response['metadata']['description']}</label>` +
                `<br><label id="resultName"><b>Author name</b>: ${response['metadata']['authors'][0]['name']}</label>` +
                `<br><label id="resultSurname"><b>Author surname</b>: ${response['metadata']['authors'][0]['surname']}</label>` +
                `<br><label id="resultDegree"><b>Author science degree</b>: ${response['metadata']['authors'][0]['scienceDegree']}</label>` +
                `<br><label id="resultUniqueness"><b>Uniqueness</b>: ${response['metadata']['uniqueness']}</label>` +
                `<br><label id="resultIpfsHash"><b>Ipfs hash</b>: <a href="${response['metadata']['ipfsHash']}">Your paper at IPFS</a></label>`
            break
        //TODO: add handling for other states
        default:
            console.log(`State: ${response['state']}. Message:  with ${response['message']}`)
            showError(response, 'addPaperResultView')
    }
}
let getNftText = document.getElementById('getNftText');
getNftText.addEventListener('keydown', function (event) {
    if (event.key === 'Enter') {
        document.getElementById('paperNftEnterView').style.display = displayStyleShow;
        document.getElementById('loadingView').style.display = displayStyleShow;
        document.getElementById('getPaperResultView').style.display = displayStyleHide;

        let http = new XMLHttpRequest();
        http.onreadystatechange = function () {
            if (http.readyState === XMLHttpRequest.DONE) {
                console.log('response status code ' + http.status)
                if (http.status !== 200 && http.status !== 202) {
                    console.log("get paper info status failed.")
                    return
                }

                handleGetPaperInfoResult(JSON.parse(http.responseText))
            }
        }

        http.open('GET',
            `${getPaperInfoEndpoint}?${paperNftKey}=${document.getElementById('getNftText').value}`, true);
        http.send(null);
        console.log('get paper info request was sent')
        getNftText.value = '';
    }
})

function handleSearchForPaperResult(response) {
    document.getElementById('searchTextEnterView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('searchResultDiv').style.display = displayStyleShow;

    switch (response['state']) {
        case okayState:
            let searchResultDiv = document.getElementById('searchForPaperResultView')
            searchResultDiv.innerHTML = `<h3>Search result for "${response['payload']}"</h3>`

            for (let i = 0; i < Math.min(response['nftMetadata'].length, response['paperMetadata'].length); i++) {
                searchResultDiv.innerHTML +=
                    `<label id="resultNftAddress"><b>Your NFT:</b> ${response['nftMetadata'][i]['address']}</label>` +
                    `<br><label id="resultNftSymbol"><b>Your NFT symbol:</b> ${response['nftMetadata'][i]['symbol']}</label>` +
                    `<br><label id="resultNftName"><b>Your genearted NFT:</b> ${response['nftMetadata'][i]['name']}</label>` +
                    `<br><label id="resultTopic"><b>Topic:</b> ${response['paperMetadata'][i]['topic']}</label>` +
                    `<br><label id="resultDescription"><b>Description:</b> ${response['paperMetadata'][i]['description']}</label>` +
                    `<br><label id="resultName"><b>Author name:</b> ${response['paperMetadata'][i]['authors'][0]['name']}</label>` +
                    `<br><label id="resultSurname"><b>Author surname:</b> ${response['paperMetadata'][i]['authors'][0]['surname']}</label>` +
                    `<br><label id="resultDegree"><b>Author science degree:</b> ${response['paperMetadata'][i]['authors'][0]['scienceDegree']}</label>` +
                    `<br><label id="resultUniqueness"><b>Uniqueness:</b> ${response['paperMetadata'][i]['uniqueness']}</label>` +
                    `<br><label id="resultIpfsHash"><b>Ipfs hash:</b> <a href="${response['paperMetadata'][i]['ipfsHash']}">Your paper at IPFS</a></label>`

                searchResultDiv.innerHTML += `<hr>`
            }
            document.getElementById('searchForPaperResultView').style.display = displayStyleShow;
            hideSearchPaperForm()
            break
        //TODO: add handling for other states
        default:
            console.log(`State: ${response['state']}. Message:  with ${response['message']}`)
            showError(response, 'addPaperResultView')
    }
}
let searchText = document.getElementById('searchText');
searchText.addEventListener('keydown', function (event) {
    if (event.key === 'Enter') {
        document.getElementById('searchTextEnterView').style.display = displayStyleShow;
        document.getElementById('loadingView').style.display = displayStyleShow;
        document.getElementById('searchForPaperResultView').style.display = displayStyleHide;

        let http = new XMLHttpRequest();
        http.onreadystatechange = function () {
            if (http.readyState === XMLHttpRequest.DONE) {
                console.log('response status code ' + http.status)
                if (http.status !== 200 && http.status !== 202) {
                    console.log("get paper info status failed.")
                    return
                }

                handleSearchForPaperResult(JSON.parse(http.responseText))
            }
        }

        http.open('GET',
            `${searchForPaperEndpoint}?${searchPayloadKey}=${document.getElementById('searchText').value}`, true);
        http.send(null);
        console.log('get paper info request was sent')
        searchText.value = '';
    }
})

function showError(response, divId) {
    document.getElementById(divId).innerHTML =
        `<h3>Ouh... Something bad happened</h3>` +
        `<label id="resultState">Your NFT: ${response['state']}</label>` +
        `<label id="resultMessage">Your NFT: ${response['resultMessage']}</label>`
}