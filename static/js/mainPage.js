function createTags(data) {
    let tags = '<div class="listOfTags"><ul>';
    tags += `<li><a data-page="filter=All">All</a></li>`;
    for (let i=0; i<data.tags.length; i++) {
        const tag = data.tags[i];
        tags += `<li><a data-page="filter=${tag.name}">${tag.name}</a></li>`
    }
    tags += '</ul></div>';
    return tags
}

function listPosts(data, rid) {
    let posts = "";

    if (data.posts === null || data.posts.length === 0) {
        return "<p>There are no posts.</p>";
    }

    for (let i=0; i<data.posts.length; i++) {
        const post = data.posts[i];
        if (rid !== "All") {
            if (!(post.tags.includes(rid))) {
                continue
            }
        }
        function postTags(post) {
            let ans = "";
            for (let j=0; j<post.tags.length; j++) {
                ans += post.tags[j];
                if (j !== post.tags.length-1) {
                    ans += " ";
                }
            }
            return ans
        }
        posts += `<div class="postPreview">
        <a data-page="viewPost=${post.id}"><h3 class="postTitle">${post.title}</h3></a>
        <p class="postPreviewInfo">
            Posted by: ${post.username} || <span class="creationDate">${post.creationDate}</span> || ${postTags(post)}
        </p>
    </div>`;
    }
    return posts
}

function mainPage(data, what) {
    const tags = createTags(data);
    const posts = '<div id="posts">' + listPosts(data, what) + '</div>';
    const newPost = '<a data-page="createPost">Create a new post</a>';
    return tags+posts+newPost
}

export { mainPage, listPosts }