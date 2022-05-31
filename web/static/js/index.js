// central containers ids
const appPaperFormContainerId = 'appPaperFormContainer'
const getPaperFormContainerId = 'getPaperFormContainer'
const searchForPaperFormContainerId = 'searchForPaperFormContainer'

// add paper containers ids
const paperUploadViewId = 'paperUploadView'
const addPaperResultViewId = "addPaperResultView"

// get paper containers ids
const paperNftEnterViewId = 'paperNftEnterView'
const getPaperResultViewId = 'getPaperResultView'

// search for paper containers ids
const searchTextEnterViewId = 'searchTextEnterView'
const searchForPaperResultId = 'searchForPaperResultView'

const loadingViewId = 'loadingView'

const displayStyleHide = 'none'
const displayStyleShow = 'block'

let authorsCnt = 1
const authorNameBaseId = 'authorName'
const authorSurnameBaseId = 'authorSurname'
const authorDegreeBaseId = 'authorDegree'

function showAppPaperForm() {
    console.log("Show add paper form")

    document.getElementById(appPaperFormContainerId).style.display = displayStyleShow;
    document.getElementById(getPaperFormContainerId).style.display = displayStyleHide;
    document.getElementById(searchForPaperFormContainerId).style.display = displayStyleHide;

    document.getElementById(paperUploadViewId).style.display = displayStyleShow;
    document.getElementById(loadingViewId).style.display = displayStyleHide;
    document.getElementById(addPaperResultViewId).style.display = displayStyleHide;

    generateAuthorFields()
}

function showGetPaperForm() {
    console.log("Show get paper form")

    document.getElementById(appPaperFormContainerId).style.display = displayStyleHide;
    document.getElementById(getPaperFormContainerId).style.display = displayStyleShow;
    document.getElementById(searchForPaperFormContainerId).style.display = displayStyleHide;

    document.getElementById(paperNftEnterViewId).style.display = displayStyleShow;
    document.getElementById(loadingViewId).style.display = displayStyleHide;
    document.getElementById(getPaperResultViewId).style.display = displayStyleHide;
}

function showSearchForPaperForm() {
    console.log("Show search for paper form")

    document.getElementById(appPaperFormContainerId).style.display = displayStyleHide;
    document.getElementById(getPaperFormContainerId).style.display = displayStyleHide;
    document.getElementById(searchForPaperFormContainerId).style.display = displayStyleShow;

    document.getElementById(searchTextEnterViewId).style.display = displayStyleShow;
    document.getElementById(loadingViewId).style.display = displayStyleHide;
    document.getElementById(searchForPaperResultId).style.display = displayStyleHide;
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
        nameDiv.innerHTML+= `<input type="text" id="${authorNameBaseId+i}" name="${authorNameBaseId+i}">`
        surnameDiv.innerHTML+= `<input type="text" id="${authorSurnameBaseId+i}" name="${authorSurnameBaseId+i}">`
        degreeDiv.innerHTML+= `<input type="text" id="${authorDegreeBaseId+i}" name="${authorDegreeBaseId+i}">`

        nameLabel.setAttribute('for', nameLabel.getAttribute('for') + `${authorNameBaseId+i} `)
        surnameDiv.setAttribute('for', surnameLabel.getAttribute('for') + `${authorSurnameBaseId+i} `)
        degreeDiv.setAttribute('for', degreeLabel.getAttribute('for') + `${authorDegreeBaseId+i} `)
    }
}

document.getElementById('uploadForm').addEventListener('submit', function (event) {
    console.log("Upload form submitted")
    event.preventDefault()
})

document.getElementById('getNftText').addEventListener('keydown', function (event) {
    console.log('Get pressed key: ' + event.key + ' Code: ' + event.code)
})

document.getElementById('searchText').addEventListener('keydown', function (event) {
    console.log('Search pressed key: ' + event.key + ' Code: ' + event.code)
})