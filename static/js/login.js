import { addNavLinks, pageLink, loadPage } from "./main.js";

const loginForm = `
    <h2>Login</h2>
    <form id="login-form">
        <label for="username">Email or nickname:</label>
        <input type="text" id="username" name="username" required><br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br>
        <button type="submit">Login</button>
    </form>
    <p class="errorMessage" id="error"></p>
`
function createLoginForm() {
    const contentDiv = document.getElementById("content");
    contentDiv.innerHTML = loginForm;

    function logSubmit(event) {
        event.preventDefault();
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        fetch(pageLink("login"), {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.error.error === true) {
                    const error = document.getElementById("error");
                    error.innerHTML = data.error.message;
                }
                if (data.user.loggedIn) {
                    loadPage("");
                }
            })
    }

    const form = document.getElementById("login-form");
    form.addEventListener("submit", logSubmit);
}

function logInHead(data) {
    const userDiv = document.getElementById("user");
    userDiv.innerHTML = `Logged in as <strong>${data.user.username}</strong>`;
    const links = document.getElementById("navBarLinks");
    links.innerHTML = '<li><a data-page="logout">Log out</a></li>';
    addNavLinks("nav li a");
}

export { createLoginForm }
export { logInHead }