import { convertTime } from "./main.js";

let wsConnection = false;
let socket;

function changeConnection(input) {
    wsConnection = input;
}

function ws(username) {
    if (!wsConnection) {

        // Store the username in a variable for later use
        socket = new WebSocket(`ws://localhost:8080/ws?username=${username}`);
        let onlineUsers = "";

        socket.addEventListener("open", (event) => {
            wsConnection = true;
        });

        socket.addEventListener('error', (error) => {
            console.error(`WebSocket error: ${error.message}`);
        });

        const onlineUsersList = document.getElementById("onlineUsersList");
        const offlineUsersList = document.getElementById("offlineUsersList");
        const onlineUsersSection = document.getElementById("onlineUsersSection");
        onlineUsersSection.style.display = "";
        const chatSection = document.getElementById("chatSection");
        const chatHeadFrom = document.getElementById("chatHeadFrom");
        const chatMessages = document.getElementById("chatMessages");
        const messageInput = document.getElementById("messageInput");
        let chattingWith = "";
        let lastLoadedMessage = 0;
        let isTyping = false;
        let typingTimer;
        let debounceTimer;

        // Stop showing typing indicator
        function stopTyping(to) {
            socket.send(JSON.stringify({ type: "stoppedTyping", user: username, to: to }));
        }

        socket.addEventListener("message", (event) => {
            if (!wsConnection) {
                socket.close();
                return
            }
            const message = JSON.parse(event.data);
            if (message.type === "onlineUsers") {
                updateOnlineUsersList(message.users);
            }

            if (message.user !== username) { // Message from not me
                if (message.user === chattingWith) { // Message from chattingWith
                    if (message.type === "message") {
                        displayMessage(message, true);
                        changeLastMessageTime(chattingWith, username, message.time, true, 0);
                    } else if (message.type === "typing") {
                        showTypingIndicator(message);
                    } else if (message.type === "stoppedTyping") {
                        hideTypingIndicator();
                    }

                } else {
                    if (message.type === "message") { // Message from someone else than chattingWith
                        const field = document.getElementById(`${message.user}Notification`);
                        const before = field.innerText;
                        let unread;
                        if (before.length === 0) {
                            unread = "1";
                        } else {
                            unread = (parseInt(field.innerText, 10) + 1).toString();
                        }
                        changeLastMessageTime(message.user, message.to, message.time, false, unread);

                    }
                }
            } else {
                if (message.type === "message") {
                    displayMessage(message, true);
                } else if (message.type === "getMessages") {
                    if (message.messages !== null) {
                        const nrOfMessages = message.messages.length;
                        lastLoadedMessage = message.messages[nrOfMessages-1].id;
                        if (message.id !== 0) {
                            const oldScrollPos = chatMessages.scrollHeight;
                            for (let i = 0; i < nrOfMessages; i++) {
                                displayMessage(message.messages[i], false);
                            }
                            chatMessages.scrollTop = chatMessages.scrollHeight - oldScrollPos;
                        } else {
                            for (let i = nrOfMessages; i > 0; i--) {
                                displayMessage(message.messages[i-1], true);
                            }
                        }
                    }
                }
            }

            function showTypingIndicator(message) {
                const typingIndicator = document.getElementById("typingIndicator");
                typingIndicator.innerHTML = `${message.user} is typing
                    <div class="dot dot1"></div>
                    <div class="dot dot2"></div>
                    <div class="dot dot3"></div>`;
            }

            function displayMessage(message, first) {
                function htmlEncode(s) {
                    const el = document.createElement("div");
                    el.innerText = el.textContent = s;
                    s = el.innerHTML;
                    return s;
                }
                const messageElement = document.createElement("div");
                let mClass = "messageOther";
                let other = message.user;
                if (message.user === username) {
                    mClass = "messageMe";
                    other = message.to;
                }
                messageElement.classList.add("message");
                messageElement.classList.add(mClass);
                messageElement.innerHTML = `<span class="messageSender">${message.user} ${convertTime(message.time)}</span>
        <div class="messageContent">${htmlEncode(message.content)}</div>`;
                if (first) {
                    chatMessages.appendChild(messageElement);
                    // Scroll the chat container to the bottom
                    chatMessages.scrollTop = chatMessages.scrollHeight;
                } else {
                    chatMessages.insertBefore(messageElement, chatMessages.firstChild);
                }
            }
        });

        function hideTypingIndicator() {
            const typingIndicator = document.getElementById("typingIndicator");
            typingIndicator.innerHTML = "";
        }

        messageInput.addEventListener("keydown", (event) => {
            const sendButton = document.getElementById("sendButton");
            sendButton.addEventListener("click", sendMessageButton);
            if (event.key === "Enter") {
                sendMessageButton();
            }

            function sendMessageButton() {
                const message = messageInput.value;

                if (message.trim() !== "") {
                    sendMessage(message);
                    messageInput.value = "";
                }
            }

            function sendMessage(message) {
                const time = new Date();
                stopTyping(chattingWith);
                const data = {
                    type: "message",
                    content: message,
                    time: time,
                    user: username,
                    to: chattingWith,
                };
                socket.send(JSON.stringify(data));
                changeLastMessageTime(chattingWith, username, time, true, 0);
            }
        });

        // Update the UI to display the list of online users
        function updateOnlineUsersList(users) {
            onlineUsersList.innerHTML = "";
            offlineUsersList.innerHTML = "";
            users.sort((a, b) => {
                if (!a.lastMessage && !b.lastMessage) {
                    return a.user.localeCompare(b.user); // Sort alphabetically if both have no last message
                } else if (!a.lastMessage) {
                    return 1; // User with no last message comes after the one with a message
                } else if (!b.lastMessage) {
                    return -1; // User with no last message comes before the one with a message
                }

                // Sort by last message time (latest messages first)
                const timeA = new Date(a.lastMessage.time);
                const timeB = new Date(b.lastMessage.time);
                return timeB - timeA;
            });
            onlineUsers = users;

            for (const user of onlineUsers) {
                if (user.user !== username) {
                    const userItem = document.createElement("li");
                    let unread = "";
                    let nClass = "";
                    if (user.lastMessage !== null && user.lastMessage.isRead === false) {
                        if (!user.lastMessage.isRead && user.lastMessage.to === username) {
                            unread = user.lastMessage.unread;
                            nClass = ' class="unreadNotifications"';
                        }
                    }
                    userItem.innerHTML = `<span>${user.user}</span> <span id="${user.user}Notification"${nClass}>${unread}</span>`;
                    if (user.online) {
                        onlineUsersList.appendChild(userItem);
                    } else {
                        offlineUsersList.appendChild(userItem);
                    }
                }
            }
        }

        function changeLastMessageTime(username, to, newTime, read, unread) {
            const userIndex = onlineUsers.findIndex((user) => user.user === username);
            let change = false;
            if (userIndex !== -1) {
                if (onlineUsers[userIndex].lastMessage) {
                    if (onlineUsers[userIndex].lastMessage.time < newTime) {
                        change = true;
                    }
                } else {
                    change = true;
                }
                if (change) {
                    onlineUsers[userIndex].lastMessage = {
                        time: newTime,
                        to: to,
                        isRead: read,
                        unread: unread,
                    }
                }
            }
            updateOnlineUsersList(onlineUsers);
        }

        // Event listener for clicking on a username in the onlineUsersList
        onlineUsersSection.addEventListener("click", (event) => {
            if (event.target.tagName === "SPAN") {
                const recipientUsername = event.target.textContent;
                if (recipientUsername !== username && recipientUsername !== chatHeadFrom.textContent) {
                    openChatWithUser(recipientUsername);
                    const field = document.getElementById(`${recipientUsername}Notification`);
                    field.innerText = "";
                    field.classList.remove("unreadNotifications");
                    messagesRead(recipientUsername, username);

                    function messagesRead(user, userTo) {
                        const data = {
                            type: "read",
                            user: user,// sender
                            to: userTo,// reader
                        };
                        socket.send(JSON.stringify(data));
                    }
                }
            }

            function loadMoreMessages() {
                clearTimeout(debounceTimer);

                debounceTimer = setTimeout(() => {
                    if (chatMessages.scrollTop === 0) {
                        getMessages(chattingWith, lastLoadedMessage);
                    }
                }, 300);
            }

            function getMessages(user, nr) {
                const data = {
                    type: "getMessages",
                    id: nr,
                    user: username,
                    to: user,
                };
                socket.send(JSON.stringify(data));
            }

            function openChatWithUser(user) {
                closeChat();
                chatMessages.removeEventListener("scroll", loadMoreMessages);
                chatHeadFrom.textContent = user;
                chatMessages.innerHTML = ""; // Clear previous messages
                chatSection.style.display = "block"; // Show the chat section
                chattingWith = user;
                getMessages(user, 0);
                chatMessages.addEventListener("scroll", loadMoreMessages);
                const field = document.getElementById(`${user}Notification`);
                field.innerText = "";
                field.classList.remove("unreadNotifications");
            }

            messageInput.addEventListener("keydown", () => {
                if (!isTyping) {
                    isTyping = true;
                    socket.send(JSON.stringify({ type: "typing", user: username, to: chattingWith }));
                } else {
                    clearTimeout(typingTimer);
                }

                // Set a timeout to send "stoppedTyping" after a period of inactivity
                typingTimer = setTimeout(() => {
                    isTyping = false;
                    stopTyping(chattingWith);
                }, 2000);
            });

            // Event listener for clicking on the chatClose button
            const chatClose = document.getElementById("chatClose");
            chatClose.addEventListener("click", () => {
                closeChat();
            });

            function closeChat() {
                chatHeadFrom.textContent = "";
                chatMessages.innerHTML = "";
                chatSection.style.display = "none";
                chattingWith = "";
                hideTypingIndicator();
                lastLoadedMessage = 0;
                messageInput.value = "";
            }
        });
    } else {
        socket.close();
        wsConnection = false;
        socket = {};
    }
}

export { ws, wsConnection, changeConnection }

