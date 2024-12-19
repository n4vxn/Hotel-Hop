# **Hotel-Hop**

**Hotel-Hop** is a Go-based web application for hotel reservations. It uses the Gin framework and PostgreSQL for efficient data management. The project offers features such as user authentication, role-based access control, and seamless handling of hotels, rooms, and bookings.

---

## **Features**

- **User Authentication**: Secure sign-up and login using JWT tokens.
- **Hotel & Room Management**: Easily add, view, and manage hotels and rooms.
- **Room Booking**: Allows users to book rooms and manage their reservations.
- **Role-Based Access Control**: Admins have extended privileges, while regular users can book and browse.
- **PostgreSQL Integration**: Reliable storage of all data.

---

## **API Routes**

### **Authentication**

- **Sign In**  
  **POST** `http://localhost:8080/api/auth/signin`  
  Log in as an existing user.

  **Request Body**:

  ```json
  {
    "email": "john@doe.com",
    "password": "password123"
  }
  ```
