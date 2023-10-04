import { createRegisterForm } from "./register.js"
import { createLoginForm, logInHead } from "./login.js"
import { mainPage, listPosts } from "./mainPage.js";
import { createPost, viewPost } from "./post.js";
import { commenting } from "./comment.js";
import { ws, wsConnection, changeConnection } from "./chat.js"

document.addEventListener("DOMContentLoaded", function () {
    // Load initial page
    addNavLinks("nav a");
    loadPage("");
});

function loadPage(pageName) {
    console.log("[Requested page] ", pageName);
    const page = pageName.split("=");
    switch (page[0]) {
        case "login":
            logOutHead();
            createLoginForm();
            break;
        case "register":
            logOutHead();
            createRegisterForm();
            break;
        case "createPost":
            createPost(storedData);
            break;
        case "filter":
            const posts = document.getElementById("posts");
            posts.innerHTML = listPosts(storedData, page[1]);
            addNavLinks("div a");
            convertDateTime();
            break;
        default:
            fetch(pageLink(pageName))
                .then(response => response.json())
                .then(data => {
                    storedData = data;
                    if (!data.httpError.error) {
                        if (data.user.loggedIn) {
                            logInHead(data);
                            const contentDiv = document.getElementById("content");
                            if (page[0] === "viewPost") {
                                const postInfo = viewPost(data)
                                contentDiv.innerHTML = postInfo;
                                commenting();
                            } else {
                                contentDiv.innerHTML = mainPage(data, "All");
                                ws(data.user.username);
                            }
                            addNavLinks("div a");
                            convertDateTime();
                        } else {
                            logOutHead();
                            if (pageName === "register") {
                                createRegisterForm();
                            } else {
                                createLoginForm();
                            }
                        }
                    } else {
                        const contentDiv = document.getElementById("content");
                        contentDiv.innerHTML = `<h2>${data.httpError.type}</h2><p>${data.httpError.text}</p><p>${data.httpError.text2}</p>`;
                    }
                })
                .catch(error => {
                    console.error("Error loading page:", error);
                });
    }
}

let storedData = ""; // Previous JSON from server
const backButton = `<a data-page="" class="backButton">&lt;&lt; Back to main page</a>`;

function pageLink(pageName) {
    return `/page?page=${pageName}`
}

function addNavLinks(selectors) {
    const navLinks = document.querySelectorAll(selectors);
    navLinks.forEach(link => {
        link.addEventListener("click", function (event) {
            event.preventDefault();
            const pageName = link.getAttribute("data-page");
            loadPage(pageName);
        });
    });
}

function logOutHead() {
    const signedOutLinks = '<li><a data-page="login">Log In</a></li>\n' +
        '        <li><a data-page="register">Register</a></li>';
    const userDiv = document.getElementById("user");
    if (userDiv.innerHTML !== "") {
        userDiv.innerHTML = "";
    }
    const links = document.getElementById("navBarLinks");
    if (links.innerHTML !== signedOutLinks) {
        links.innerHTML = signedOutLinks;
        addNavLinks("nav li a");
    }

    createLoginForm();
    if (wsConnection) {
        ws();
    }
    const onlineUsersList = document.getElementById("onlineUsersSection");
    onlineUsersList.style.display = "none";
    const chatSection = document.getElementById("chatSection");
    chatSection.style.display = "none";

}

// Function to convert ISO 8601 date-time format to a readable format
function convertDateTime() {
    const dateElements = document.querySelectorAll('.creationDate');
    dateElements.forEach(function (element) {
        const isoDateTime = element.textContent.trim();
        element.textContent = convertTime(isoDateTime);
    });
}
function convertTime(dateTimeString) {
    const dateTime = new Date(dateTimeString);
    const options = { day: 'numeric', month: 'short', year: 'numeric', hour: 'numeric', minute: 'numeric' };
    return dateTime.toLocaleDateString('en-US', options);
}

export { addNavLinks, pageLink, loadPage, backButton, convertDateTime, convertTime };


