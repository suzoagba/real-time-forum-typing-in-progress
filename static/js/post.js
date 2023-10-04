import { loadPage, pageLink, backButton, addNavLinks } from "./main.js";
import { comments } from "./comment.js";

const createPostForm = (data) => {
    function tags(data) {
        let ans = "";
        for (let i=0; i<data.tags.length; i++) {
            const tag = data.tags[i];
            ans += `<input type="checkbox" name="tags" value="${tag.id}"> ${tag.name}<br>`
        }
        return ans
    }
    return `${backButton}
    <h2>Create a New Post</h2>
    <div class="createPost">
        <form id="createPost-form">
            <label for="title">Title:</label>
            <br>
            <input type="text" id="title" name="title" required>
            <br>
            <label for="description">Description:</label>
            <br>
            <textarea id="description" name="description" maxlength="300" required></textarea>
            <br>
            <label>Tags:</label>
            <br>
            ${tags(data)}
            <button type="submit">Create Post</button>
        </form>
        <p class="errorMessage" id="error"></p>
    </div>`
}

function createPost(data) {
    const contentDiv = document.getElementById("content");
    contentDiv.innerHTML = createPostForm(data);
    addNavLinks("div a");

    function logSubmit(event) {
        event.preventDefault();
        if (validateForm()) {
            const title = document.getElementById("title").value;
            const description = document.getElementById("description").value;
            const tags = getCheckedTags("tags");

            fetch(pageLink("createPost"), {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ title, description, tags }),
            })
                .then(response => response.json())
                .then(data => {

                    if (data.error.error === true) {
                        const error = document.getElementById("error");
                        error.innerHTML = data.error.message;
                    } else if (data.user.loggedIn) {
                        loadPage("");
                    }
                })
        }
    }

    const form = document.getElementById("createPost-form");
    form.addEventListener("submit", logSubmit);
}

function viewPost(data) {
    let ans = backButton;
    const post = data.posts[0];
    ans += `<h2>${post.title}</h2>`;
    ans += `<div class="post">
                <div class="commentInfo">
                      <p class="postInfo">
                        <span class="author">${post.username}</span>
                        <br>
                        <span class="creationDate">${post.creationDate}</span>
                      </p>
                      <p class="description">
                        ${post.description}
                      </p>
                </div>
            </div>`;
    ans += comments(data);
    return ans
}


// Function to check if at least one of the tags is selected when creating a new post
function validateForm() {
    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    let checkedCount = 0;
    for (let i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            checkedCount++;
        }
    }
    if (checkedCount === 0) {
        alert("Please select at least one tag.");
        return false; // Prevent form submission
    }
    return true; // Allow form submission
}

function getCheckedTags(name) {
    const checkboxes = document.getElementsByName(name);
    const checkedTags = [];

    for (let i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            checkedTags.push(checkboxes[i].value);
        }
    }

    return checkedTags;
}

export { createPost, viewPost }