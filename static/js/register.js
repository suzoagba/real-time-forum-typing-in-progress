import { logInHead } from "./login.js";
import {loadPage, pageLink} from "./main.js";

const registerForm = `
    <h2>Register</h2>
    <form id="register-form">
        <label for="nickname">Nickname (username):</label>
        <input type="text" id="nickname" name="nickname" maxlength="10" required><br>
        <label for="age">Age:</label>
        <input type="number" id="age" name="age" min="13" max="120" required><br>
        <label for="gender">Gender:</label>
        <select id="gender" name="gender" required>
            <option disabled selected value></option>
            <option value="female">Female</option>
            <option value="male">Male</option>
            <option value="other">Other</option>
        </select><br>
        <label for="first_name">First Name:</label>
        <input type="text" id="first_name" name="first_name" required><br>
        <label for="last_name">Last Name:</label>
        <input type="text" id="last_name" name="last_name" required><br>
        <label for="email">Email:</label>
        <input type="email" id="email" name="email" required><br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br>
        <button type="submit">Register</button>
    </form>
    <p class="errorMessage" id="error"></p>
`
function createRegisterForm() {
    const contentDiv = document.getElementById("content");
    contentDiv.innerHTML = registerForm;

    function logSubmit(event) {
        //console.log("login form submit");
        event.preventDefault();
        const username = document.getElementById("nickname").value;
        const age = document.getElementById("age").value;
        const gender = document.getElementById("gender").value;
        const firstName = document.getElementById("first_name").value;
        const lastName = document.getElementById("last_name").value;
        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        fetch(pageLink("register"), {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, age, gender, firstName, lastName, email, password }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.error.error === true) {
                    const error = document.getElementById("error");
                    error.innerHTML = data.error.message;
                }
                console.log(data)
                if (data.user.loggedIn) {
/*                    const contentDiv = document.getElementById("content");
                    contentDiv.innerHTML = "<p>Success!</p>";
                    logInHead(data);*/
                    loadPage("");
                }
            })
    }

    const form = document.getElementById("register-form");
    form.addEventListener("submit", logSubmit);
}

export { createRegisterForm }