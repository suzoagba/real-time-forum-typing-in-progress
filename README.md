# real-time-forum

This project consists in creating a web forum that allows:

- Registration and Login
- Creation of posts
- Commenting posts
- Private Messages

### Allowed Packages

- SQLite
- Golang
  - All standard go packages are allowed
  - Gorilla websocket
  - sqlite3
  - bcrypt
  - UUID
- Javascript
- HTML
- CSS

### Usage

- `go run .`
- `http://localhost:8080/`
- The user is able to send private messages to the users who are online.

#### Registered users
- Username: `admin` Password: `admin`
- Username: `oper` Password: `oper` (over 50 private messages with admin)
- Username: `user` Password: `psw`
- Username: `auto` Password: `psw`
- Username: `jaan.tamm` Password: `psw`

### Audit

Questions can be found:
- [real-time-forum](https://github.com/01-edu/public/blob/master/subjects/real-time-forum/audit.md)
- [typing-in-progress](https://github.com/01-edu/public/blob/master/subjects/real-time-forum/typing-in-progress/audit.md)

## Developers
- Willem Kuningas / *thinkpad*
- Samuel Uzoagba / *suzoagba*