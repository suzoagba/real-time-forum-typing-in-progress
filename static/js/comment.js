import { loadPage, pageLink, convertDateTime } from "./main.js";

function comments(data) {
    return `<div id="comments">${loadComments(data)}</div>
<div id="reply">${addReplyOption(data.posts[0].id)}</div>`
}

function loadComments(data) {
    let ans = "";
    for (let i=0; i<data.comments.length; i++) {
        const comment = data.comments[i];
        ans += `<div class="commentInfo" id="comment${comment.id}">
                    <p class="postInfo">
                      <span class="author">${comment.username}</span>
                      <br>
                      <span class="creationDate">${comment.creationDate}</span>
                    </p>
                    <p class="commentText">
                      ${comment.content}
                    </p>
                </div>`;
    }
    return ans
}

function addReplyOption(id) {
    return `<h3>Reply:</h3>
            <form id="comment-form">
              <input type="hidden" name="postID" id="postID" value="${id}">
              <textarea name="content" id="comment" placeholder="Enter your reply" maxlength="300" required></textarea>
              <br>
              <button type="submit">Submit</button>
            </form>
            <p class="errorMessage" id="error"></p>`;
}

function commenting() {
    function commentSubmit(event) {
        event.preventDefault();
        const comment = document.getElementById("comment").value;
        const id = document.getElementById("postID").value;

        fetch(pageLink("comment"), {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ comment, id }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.error.error === true) {
                    const error = document.getElementById("error");
                    error.innerHTML = data.error.message;
                }
                console.log(data)
                if (data.user.loggedIn) {
                    const comments = document.getElementById("comments");
                    comments.innerHTML = loadComments(data);
                    convertDateTime();
                    document.getElementById("comment").value = "";
                }
            })
    }

    const form = document.getElementById("comment-form");
    form.addEventListener("submit", commentSubmit);
}

export { comments, commenting }