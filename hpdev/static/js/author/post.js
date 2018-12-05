class UI {
    constructor(document){
        this.dom = {}

        this.dom.postID = document.getElementById('post-id')
        this.dom.postTitle = document.getElementById('post-title')
        this.dom.postContent = document.getElementById('post-content')
        this.dom.convertBtn = document.getElementById('convertHTMLBtn')
        this.dom.renderingArea = document.getElementById('renderingArea')
    }

    setPostID(id) {
        this.dom.postID.value = id
    }
    setPostTitle(title) {
        this.dom.postTitle.value = title
    }
    setPostContent(content) {
        this.dom.postContent.value = content
    }

    markdown() {
        return this.dom.postContent.value
    }
    renderPostHTML(html) {
        this.dom.renderingArea.textContent = html
    }

}

const addEvents = ui => {
    ui.dom.convertBtn.addEventListener('click', convertHTML(ui), false)
}

const convertHTML = ui => {
    return () => {
        const md = ui.markdown()

        fetch('/markdown?dst=html', {
            method: 'POST',
            headers: {
                "Content-Type": "application/octet-stream",
            },
            body: md,
        })
        .then(response => {
            if (!response.ok) {
                console.log(response)
                return
            }

            response.text().then(html => {
                ui.renderPostHTML(html)
            })
        })
        .catch(error => {
            console.log(error)
        })
    }
}

const init = () => {
    const ui = new UI(document)

    addEvents(ui)
}

init()