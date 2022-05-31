const displayStyleHide = 'none'
const displayStyleShow = 'block'

let authorsCnt = 1
const authorNameBaseId = 'authorName'
const authorSurnameBaseId = 'authorSurname'
const authorDegreeBaseId = 'authorDegree'

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

function showGetPaperForm() {
    console.log("Show get paper form")

    document.getElementById('appPaperFormContainer').style.display = displayStyleHide;
    document.getElementById('getPaperFormContainer').style.display = displayStyleShow;
    document.getElementById('searchForPaperFormContainer').style.display = displayStyleHide;

    document.getElementById('paperNftEnterView').style.display = displayStyleShow;
    document.getElementById('loadingView').style.display = displayStyleHide;
    document.getElementById('getPaperResultView').style.display = displayStyleHide;
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