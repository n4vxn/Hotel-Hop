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
### Get All Hotels
**GET** `http://localhost:8080/api/v1/hotel`

### Get Hotel by ID
**GET** `http://localhost:8080/api/v1/hotel/:id`

### Get Rooms of a Hotel
**GET** `http://localhost:8080/api/v1/hotel/:id/rooms`

## Room Management

### Get All Available Rooms
**GET** `http://localhost:8080/api/v1/room`

### Book a Room
**POST** `http://localhost:8080/api/v1/room/:id/book`  
**Requires user authentication.**

### Book a Room Request Body
```json
{
  "fromDate": "2024-12-12T00:00:00Z",
  "tillDate": "2024-12-25T00:00:00Z",
  "numPersons": 4
}
```
## View Bookings by User
**GET** `http://localhost:8080/api/v1/booking/:id`

### Cancel Booking
**GET** `http://localhost:8080/api/v1/booking/:id/cancel`

## Admin Panel

### Get All Bookings
**GET** `http://localhost:8080/api/v1/admin/booking`


## Setup

### Clone the Repository
```bash
git clone https://github.com/your-repo/hotel-hop.git
cd hotel-hop

### Install Dependencies
bash
go mod tidy
```