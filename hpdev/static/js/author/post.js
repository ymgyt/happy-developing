class Tag {
    constructor(tag){
        this.data = tag
        this.data.selected = false

        this.dom = document.createElement('button')
        this.dom.textContent = this.data.name
        this.dom.setAttribute('type', 'button')
        this.dom.setAttribute('selected', this.data.selected)

        this.myOnClick = this.myOnClick.bind(this) // https://qiita.com/tsin1rou/items/90576b6c00b895478610#class
        this.dom.addEventListener('click', this.myOnClick,false)
    }

    myOnClick(event){ this.toggle(!this.data.selected) }
    enable() { this.toggle(true) }
    disable() { this.toggle(false) }
    toggle(flag) {
        this.data.selected = flag
        this.dom.setAttribute('selected', flag)
    }
    getDom() { return this.dom }
}

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
        this.dom.createTag = document.getElementById('create-tag')
        this.dom.tagSelect = document.getElementById('tag-select')

        this.dom.savePostBtn = document.getElementById('savePostBtn')
        this.dom.createTagBtn = document.getElementById('createTagBtn')

        this.tags = []
        this.postTagIDs = []
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
    setPostContent(content) {
        this.dom.postContent.value = content
    }
    renderPostHTML(html) {
        this.dom.renderingArea.innerHTML = html
    }
    setTags(tags) {
        for (const tagData of tags) {
            const tag = new Tag(tagData)
            this.dom.tagSelect.appendChild(tag.getDom())
            this.tags.push(tag)
        }
        this.enableTags(this.postTagIDs)
    }
    enableTags(tagIDs = []) {
        this.postTagIDs = tagIDs
        this.tags.forEach(tag => {
            if (tagIDs.includes(tag.data.id)) {
                tag.enable()
            } else {
                tag.disable()
            }
        })
    }
    clearTags() {
        this.tags = []
        this.dom.tagSelect.innerHTML = ''
    }

    postID() { return Number(this.dom.postID.value) }
    postTitle() { return this.dom.postTitle.value }
    postURLSafeTitle() { return this.dom.postURLSafeTitle.value }
    postCreatedAt() { return this.dom.postCreatedAt.value }
    postStatus() {
        const opts = this.dom.postStatus.options
        for(let i=0;i<opts.length;i++) {
            if (opts[i].selected === true) {
                return opts[i].value
            }
        }
    }
    postSelectedTagIDs() {
        return this.tags.filter(tag => tag.data.selected).map(tag => tag.data.id)
    }
    tagName() { return this.dom.createTag.value }
    clearTagName() { this.dom.createTag.value = '' }

    markdown() { return this.dom.postContent.value }
    html() { return this.dom.renderingArea.innerHTML }


    setPost(post) {
        this.setPostID(post.meta.id)
        this.setPostTitle(post.meta.title)
        this.setPostURLSafeTitle(post.meta.url_safe_title)
        this.setPostCreatedAt(post.meta.created_at)
        this.setPostStatus(post.meta.status)
        this.setPostContent(post.content.markdown)
        this.renderPostHTML(post.content.html)
        this.enableTags(post.meta.tag_ids || []) // 関数側でケアしたいが、nullだとdefault argが適用されない.
    }


}

const addEvents = ui => {
    ui.dom.savePostBtn.addEventListener('click', savePost(ui), false)
    ui.dom.createTagBtn.addEventListener('click', saveTag(ui), false)
    ui.dom.postContent.addEventListener('keydown', textareaHandler(ui), false)
}

const tagBtnOnClick = (tag,tagBtn) => {
    return () => {
        tag.selected = !tag.selected
        tagBtn.setAttribute('selected', tag.selected)
    }
}

const isCreate = () => {
    return window.location.pathname.endsWith("new")
}

const savePost = ui => {
    return () => {
        const ep = "/api" + window.location.pathname
        let post = {}
        // mergeとか利用して、共通化したい
        if (isCreate()) {
            post = {
                "meta": {
                    "id": 0,
                    "title": ui.postTitle(),
                    "url_safe_title": ui.postURLSafeTitle(),
                    "created_at": new Date(), // TODO JST固定
                    "status": ui.postStatus(),
                    "tag_ids": ui.postSelectedTagIDs(),
                },
                "content": {
                    "markdown": ui.markdown(),
                    "html": ui.html(),
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
                    "tag_ids": ui.postSelectedTagIDs(),
                },
                "content": {
                    "markdown": ui.markdown(),
                    "html": ui.html(),
                },
            }
        }

        console.log("save post", post)

        fetch(ep, {
            method: isCreate() ? 'POST' : 'PUT',
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
                ui.setPost(post)
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
            ui.setPost(post)
        })
    })
    .catch(error => console.log(error))
}

const saveTag = ui => {
    return () => {
        const ep = "/api/author/tags"
        const tag = {
            "name": ui.tagName(),
        }
        if (tag.name === "") {
            console.log("empty tag name")
            return
        }
        fetch(ep, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(tag)
        })
        .then(response => {
            if (!response.ok) {
                response.json().then(json => console.log("fail to create", json))
                return
            }
            ui.clearTagName()
            fetchTags(ui)
        })
        .catch(error => console.log(error))
    }
}

const fetchTags = ui => {
    const ep = "/api/author/tags"
    fetch(ep)
    .then(response => {
        if (!response.ok) {
            console.log(response)
            return
        }
        response.json().then(tags => {
            tags.sort((x,y) => x.name.localeCompare(y.name))
            ui.clearTags()
            ui.setTags(tags)
        })
    })
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

const textareaHandler = ui => {
    const activate = convertHTML(ui)
    return event => {
        // hot convert
        if (event.key === "Enter" && event.shiftKey) {
            activate()
        }

        // tab insertion
        if (event.key === "Tab") {
            event.preventDefault()
            const t = event.target
            const start = t.selectionStart
            const end = t.selectionEnd
            t.value = t.value.substring(0,start) + "\t" + t.value.substring(end)
            t.selectionEnd = start + 1
        }
    }
}



let gUI ={} // for debug

const init = () => {
    const ui = new UI(document)
    gUI = ui

    // 新規作成の場合はcallしない
    if (!isCreate()) {
        fetchPost(ui)
    }
    fetchTags(ui)
    addEvents(ui)
}

init()