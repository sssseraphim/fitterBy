# **GymBro API** - Fitness Social Network

A RESTful API for fitness enthusiasts to create programs, track workouts, and connect with others.

---

## 🌐 **Base URL**
```
http://localhost:8080/api
```

All endpoints return JSON responses.

---

## **Authentication**

### **Sign Up**
Create a new user account.

```http
POST /auth/signup
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}
```

### **Login**
Authenticate and receive JWT token.

```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
    "id": "776ceb50-95b6-4e70-9b92-d21858a25fc6",
    "name": "John Doe",
    "email": "user@example.com",
    "TokenResoponse": {
       "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzc2Y2ViNTAtOTViNi00ZTcwLTliOTItZDIxODU4YTI1ZmM2IiwidXNlcl90eXBlIjoidXNlcl90eXBlIiwiZW1haWwiOiJlbWFpbDEiLCJzdWIiOiI3NzZjZWI1MC05NWI2LTRlNzAtOWI5Mi1kMjE4NThhMjVmYzYiLCJleHAiOjE3NjQ4MTA4NTksImlhdCI6MTc2NDgwNzI1OSwianRpIjoiOTk2NDRmMDMtZDdjZi00OGZkLWFlNGItY2Q1OTIwODQwNDU2In0.m4IMZEjSx8vbGu4WFS87rc4PekIHdLQUTpQR463EIo8",
        "token_type": "Bearer",
        "expires_in": 3600
    }
}
```

**Note:** Include token in headers for protected endpoints:
```
Authorization: Bearer <your_token>
```

---

## 👤 **User Endpoints**

### **Get User Profile**
```http
GET /users/{user_id}
```
Get public profile of any user.

### **Get Current User**
```http
GET /me
```
**Protected** - Get authenticated user's own profile.

### **Update Bio**
```http
PATCH /me/bio
```
**Protected** - Update your biography.

**Request Body:**
```json
{
  "bio": "Just a guy trying to get huge "
}
```

### **Follow User**
```http
POST /users/follow
```
**Protected** - Follow another user.

**Request Body:**
```json
{
  "user_id": "uuid-of-user-to-follow"
}
```

### **Get Followed Users**
```http
GET /users/follow
```
**Protected** - Get list of users you're following.

---

## **Post Endpoints**

### **Get Single Post**
```http
GET /posts/{post_id}
```
Get details of a specific post.

### **Get Followed Posts**
```http
GET /posts/followed
```
**Protected** - Get posts from users you follow.

### **Create Post**
```http
POST /posts
```
**Protected** - Create a new post.

**Request Body:**
```json
{
  "content": "Just hit a new PR! 225lbs bench! 🏋🏾‍♂️",
  "media_urls": ["https://example.com/pr.jpg"]
}
```

### **Like Post**
```http
POST /posts/like
```
**Protected** - Like/unlike a post.

**Request Body:**
```json
{
  "post_id": "uuid-of-post"
}
```

### **Comment on Post**
```http
POST /posts/comments
```
**Protected** - Add a comment to a post.

**Request Body:**
```json
{
  "post_id": "uuid-of-post",
  "content": "Great work bro!"
}
```

### **Get Post Comments**
```http
GET /posts/comments
```
Get comments for a specific post.

**Request Body:**
```json
{
  "post_id": "uuid-of-post",
}
```
---

##  **Exercise Endpoints**

### **Create Exercise**
```http
POST /exercises
```
**Protected** - Create a new exercise in the database.

**Request Body:**
```json
{
  "name": "Barbell Bench Press",
  "description": "Flat bench press with barbell",
  "media_urls": "https://example.com/bench-press.mp4"
}
```

### **Get All Exercises**
```http
GET /exercises
```
Get list of all exercises in the database.


### **Get Exercise by ID**
```http
GET /exercises/{exercise_id}
```
Get details of a specific exercise.

---

## **Program Endpoints**

### **Create Program**
```http
POST /programs
```
**Protected** - Create a workout program.

**Request Body:**
```json
{
  "name": "Big Boobs Program",
  "description": "Get that chest pump going!",
  "visibility": "public",
  "days": [
    {
      "name": "Chest Day",
      "description": "Bench and flyes",
      "order": 1,
      "lifts": [
        {
          "exercise_id": "uuid-of-bench-press",
          "description": "Heavy bench 5x5",
          "sets": 5,
          "reps": 5,
          "order": 1
        }
      ]
    }
  ]
}
```

### **Get All Programs**
```http
GET /programs
```

### **Get Program by ID**
```http
GET /programs/{program_id}
```
Get details of a specific program.

### **Subscribe to Program**
```http
POST /programs/{program_id}/subscribe
```
**Protected** - Start following/doing a program.

### **Get My Subscribed Programs**
```http
GET /users/me/programs
```
**Protected** - Get programs you're currently following.

---

## **Workout Endpoints**

### **Create Workout**
```http
POST /workouts
```
**Protected** - Log a completed workout.

**Request Body:**
```json
{
    "program_id": "c5c25979-139d-4bce-a27b-5f1657418771",
    "program_day_id": "bb07f054-fe20-4969-b4fd-4e8f8a1c96ff",
    "lifts":[
        {
            "exercise_id": "4471ddaf-3901-4572-8479-f0ecd39ecf70",
            "weight": 100,
            "sets": 5,
            "reps": 5
        }
    ]
}
```

### **Get My Workouts**
```http
GET /users/me/workouts
```
**Protected** - Get your workout history.


### **Get Workout by ID**
```http
GET /workouts/{workout_id}
```
**Protected** - Get details of a specific workout.


