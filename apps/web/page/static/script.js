const state = {}
const formValueTraversal = ["inputText.value", "outputText.value", "key.value"]
let isFileInput = false
let selectedTab = "corrupt"

document.addEventListener("DOMContentLoaded", () => {
    const inputForm = document.querySelector("form#form")
    const tabButtons = document.querySelectorAll(".tabs>button.tab-button")
    const animatedPageTitle = document.querySelector(".page-title>span")
    let pageTitleAnimationInterval;

    function togglePageTitleAnimation(tab) {
        if(tab === "uncorrupt") {
            clearInterval(pageTitleAnimationInterval);
            animatedPageTitle.style.opacity = 1
            pageTitleAnimationInterval= null
        } else {
            pageTitleAnimationInterval = setInterval(() => {
                animatedPageTitle.style.opacity = Math.random()
            }, [50])
        }
    }
    togglePageTitleAnimation(selectedTab)

    tabButtons.forEach(btn => {
        btn.addEventListener("click", (ev) => {
            let target = ev.target.getAttribute("data-target"), current = null;
            if(selectedTab === target) return
            selectedTab = target
            togglePageTitleAnimation(selectedTab)
            
            tabButtons.forEach(b => {
                if(b.getAttribute("data-target") === target) {
                    b.classList.add("active")
                } else {
                    if(b.classList.contains("active")) {
                        current = b.getAttribute("data-target")
                    }
                    b.classList.remove("active")
                }
            })

            if(current){
                let newCurrentState = {}
                for(const v of formValueTraversal) {
                    const paths = v.split('.')
                    let obj = {}
                    let currentItem = form, nestedObj = obj, nextState = state[target];
                    for(let i=0; i<paths.length-1; i++) {
                        currentItem = currentItem[paths[i]]
                        nestedObj[paths[i]] = {}
                        nestedObj = nestedObj[paths[i]]
                        nextState = nextState?.[paths[i]]
                    }
                    nestedObj[paths[paths.length-1]] = currentItem[paths[paths.length-1]]
                    currentItem[paths[paths.length-1]] = nextState?.[paths[paths.length-1]] || ""
                    newCurrentState = { ...newCurrentState, ...obj }
                }
                state[current] = { ...newCurrentState }
            }
            inputForm.setAttribute("data-target", target)
        })
    })
    
    document.querySelector(".input-type-selector>input").addEventListener("change", (ev) => {
        isFileInput = ev.target.checked
        if(ev.target.checked) {
            document.querySelector('.input-type.input-text').classList.remove('active')
            document.querySelector('.input-type.input-file').classList.add('active')
        }
        else {
            document.querySelector('.input-type.input-file').classList.remove('active')
            document.querySelector('.input-type.input-text').classList.add('active')
        }
    })
    inputForm.addEventListener("submit", (ev) => {
        ev.preventDefault()
        const file = ev.target.inputFile?.files?.[0]
        const text = ev.target?.inputText?.value
        if((isFileInput && !file) || (!isFileInput && !text)) {
            alert('No valid input provided')
            return
        }
        const formData = new FormData()
        formData.append('key', ev.target.key?.value)
        if(isFileInput) {
            formData.append('inputFile', file)
        }
        else {
            formData.append('inputText', text)
        }
        
        fetch(`/${inputForm.getAttribute('data-target')}`, {
            method: 'POST',
            body: formData
        })
        .then(res =>  {
            if(isFileInput) return res.arrayBuffer()
            else return res.text()
        })
        .then(res => {
            if(isFileInput) {
                handleDownload(res, file.name)
            }
            else {
                ev.target.outputText.value = res
            }
        })
        .catch(err => {
                alert(err)
        })
    })
    function handleDownload(data, filename='example') {
        const blob = new Blob([data], { type: 'application/octet-stream'})
        const url = URL.createObjectURL(blob)
        const downloadLink = document.createElement('a');
        downloadLink.href = url;
        downloadLink.download = filename
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
        URL.revokeObjectURL(url);
    }
})