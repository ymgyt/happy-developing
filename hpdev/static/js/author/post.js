class UI {
    constructor(document){
        this.dom = {}

        this.dom.postID = document.getElementById('post-id')
        this.dom.postTitle = document.getElementById('post-title')
        this.dom.postURLSafeTitle = document.getElementById('post-url-safe-title')
        this.dom.postCreatedAt = document.getElementById('post-created-at')
        this.dom.postStatus = document.getElementById('post-status')
        this.dom.postContent = document.getElementById('post-content')
        this.dom.renderingArea = document.getElementById('renderingArea')

        this.dom.convertBtn = document.getElementById('convertHTMLBtn')
        this.dom.savePostBtn = document.getElementById('savePostBtn')
    }

    setPostID(id) {
        this.dom.postID.value = id
    }
    setPostTitle(title) {
        this.dom.postTitle.value = title
    }
    setPostURLSafeTitle(safeTitle) {
        this.dom.postURLSafeTitle.value = safeTitle
    }
    setPostCreatedAt(at) {
        this.dom.postCreatedAt.value = at
    }
    setPostContent(content) {
        this.dom.postContent.value = content
    }
    setPostStatus(status) {
        if (status === "undefined") {
            status = "draft"
        }
        const opts = this.dom.postStatus.options
        for(let i=0;i<opts.length;i++) {
            if (opts[i].value === status) {
                opts[i].selected = true
            }
        }
    }

    postID() { return this.dom.postID.value }
    postTitle() { return this.dom.postTitle.value }
    postURLSafeTitle() { return this.dom.postURLSafeTitle.value }
    postCreatedAt() { return this.dom.postCreatedAt.value }
    postContent() { return this.dom.postContent.value }
    postStatus() {
        const opts = this.dom.postStatus.options
        for(let i=0;i<opts.length;i++) {
            if (opts[i].selected === true) {
                return opts[i].value
            }
        }
    }

    markdown() {
        return this.dom.postContent.value
    }
    renderPostHTML(html) {
        this.dom.renderingArea.innerHTML = html
    }

}

const addEvents = ui => {
    ui.dom.convertBtn.addEventListener('click', convertHTML(ui), false)
    ui.dom.savePostBtn.addEventListener('click', savePost(ui), false)
}

const savePost = ui => {
    return () => {
        const ep = "/api" + window.location.pathname
        const isCreate = ep.endsWith("new")
        let post = {}
        if (isCreate) {
            post = {
                "meta": {
                    "id": 0,
                    "title": ui.postTitle(),
                    "url_safe_title": ui.postURLSafeTitle(),
                    "created_at": new Date(), // TODO JST固定
                    "status": ui.postStatus(),
                },
                "content": {
                    "markdown": ui.postContent(),
                },
            }
        } else {
            post = {
                "meta": {
                    "id": ui.postID(),
                    "title": ui.postTitle(),
                    "url_safe_title": ui.postURLSafeTitle(),
                    "created_at": ui.postCreatedAt(),
                    "status": ui.postStatus(),
                },
                "content": {
                    "markdown": ui.postContent(),
                },
            }
        }

        console.log("save post", post)

        fetch(ep, {
            method: isCreate ? 'POST' : 'PUT',
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(post),
        })
        .then(response => {
            response.json().then(post => {
                if (response.ok && isCreate) {
                    window.history.pushState({},'', `/author/posts/${post.meta.id}`)
                }
                console.log(post)
            })
        })
        .catch(error => console.log(error))
    }
}

const fetchPost = ui => {
    const ep = "/api" + window.location.pathname
    fetch(ep)
    .then(response => {
        if (!response.ok) {
            console.log(response)
            return
        }
        response.json().then(post => {
            ui.setPostID(post.meta.id)
            ui.setPostTitle(post.meta.title)
            ui.setPostURLSafeTitle(post.meta.url_safe_title)
            ui.setPostCreatedAt(post.meta.created_at)
            ui.setPostStatus(post.meta.status)
            ui.setPostContent(post.content.markdown)
            ui.renderPostHTML(post.content.html)
        })
    })
    .catch(error => console.log(error))
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

    // 新規作成の場合はcallしない
    if (!window.location.pathname.endsWith("new")) {
        fetchPost(ui)
    }
    addEvents(ui)
}

init()