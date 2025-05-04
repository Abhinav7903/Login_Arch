# 🚀 Welcome to LoginArch

**LoginArch** is a lightweight user authentication service that lets you **create users**, **log in**, **log out**, and **fetch user details** with ease.

It’s designed to be clean, quick, and ready for integration into your web or mobile projects.

---

## ✨ Features

- ✅ Create new user accounts
- 🔐 Log in with secure credentials
- 📥 Fetch user details (with proper date/time formatting)
- 🔓 Log out to end a session

---

## 🔍 How It Works (Simple Overview)

LoginArch handles users and their sessions:

- When someone **signs up**, their password is stored securely.
- When someone **logs in**, a session is created to keep them logged in.
- When someone **logs out**, the session ends.
- When someone checks their **profile**, they get user info with readable timestamps.

It ensures only valid users can log in and avoids duplicate active sessions.

---

## ⚙️ What’s Used

LoginArch is built using a few powerful tools:

- **Routing & HTTP** – Handles incoming requests and returns clear responses.
- **Session Store** – Uses a fast in-memory system to track who is logged in (Redis).
- **Data Storage** – Stores user information securely in a PostgreSQL database.
- **Structured Logging** – Keeps logs clean and readable for monitoring and debugging.

> Don’t worry — it’s already wired up. You just run it and it works.

---

## 🚦 How to Set Up

### 1. Set Up PostgreSQL (psql)

To run this project, you will need a PostgreSQL database.

#### Install PostgreSQL

- **macOS**: `brew install postgresql`
- **Linux**: Follow instructions [here](https://www.postgresql.org/download/linux/)
- **Windows**: Download from [PostgreSQL official site](https://www.postgresql.org/download/windows/)

#### Create a Database

1. After installing PostgreSQL, run the following command to open the PostgreSQL prompt:
   ```bash
   psql -U postgres

# 🛠 Install Redis also

Redis will be used for managing sessions in the LoginArch application.



### Future Improvements
![alt text](image.png)