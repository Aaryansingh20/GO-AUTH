# 🔐 Go + Gin + MongoDB – JWT Authentication API

A complete authentication backend built using **Go (Golang)**, **Gin**, **MongoDB**, **JWT**, and **middleware-based token validation**.

This project includes:

- User Signup  
- User Login  
- JWT Access Token + Refresh Token  
- Protected Routes  
- Role-Based User Types  
- Get All Users  
- Get Single User by ID  
- MongoDB database integration  

Perfect for learning backend development, using in real-world applications, or integrating with a React/Next.js frontend.

---

## 🚀 Features

### ✔ User Signup  
- Validates input  
- Hashes password using bcrypt  
- Stores user in MongoDB  
- Generates JWT + refresh token

### ✔ User Login  
- Verifies email/password  
- Generates new JWT  
- Returns token + refresh_token + user details

### ✔ JWT Authentication Middleware  
- Checks for `Authorization: Bearer <token>`  
- Validates signature  
- Extracts user claims (email, userType, uid, etc.)  
- Rejects unauthorized requests

### ✔ Protected Routes  
- `GET /users` → Get all users (Admin access recommended)  
- `GET /users/:id` → Get single user by user_id  

### ✔ MongoDB Integration  
- Uses official Go Mongo driver  
- Stores users in a `user` collection  

---

## 📂 Project Structure

go-auth/
│
├── controllers/
│ └── userController.go
│
├── database/
│ └── databaseConnection.go
│
├── helpers/
│ ├── authHelper.go
│ └── tokenHelper.go
│
├── middleware/
│ └── authMiddleware.go
│
├── models/
│ └── userModel.go
│
├── routes/
│ ├── authRoutes.go
│ └── userRoutes.go
│
├── main.go
├── go.mod
├── go.sum
└── .env (NOT committed to GitHub)
