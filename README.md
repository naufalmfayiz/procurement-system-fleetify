# Procurement System (Sistem Pengadaan Barang)

Aplikasi web sederhana untuk mencatat pembelian barang (Procurement) dari Supplier.

## 📁 Project Structure

```
procurement-system-fleetify/
├── backend/                 # Go Fiber API
│   ├── config/             # Configuration
│   ├── database/           # Database connection
│   ├── handlers/           # Request handlers
│   ├── middleware/         # Auth middleware
│   ├── models/             # GORM models
│   ├── routes/             # API routes
│   ├── main.go             # Entry point
│   ├── go.mod              # Go modules
│   └── .env.example        # Environment template
│
├── frontend/               # jQuery + Bootstrap UI
│   ├── css/               # Stylesheets
│   ├── js/                # JavaScript files
│   ├── index.html         # Login page
│   ├── register.html      # Register page
│   ├── dashboard.html     # Dashboard
│   ├── items.html         # Items CRUD
│   ├── suppliers.html     # Suppliers CRUD
│   ├── purchase.html      # Create purchase
│   └── history.html       # Purchase history
│
└── README.md
```

## 🛠️ Tech Stack

### Backend

- **Language**: Go (Golang)
- **Framework**: Go Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Token)

### Frontend

- **Library**: jQuery 3.7.1
- **Styling**: Bootstrap 5.3.2
- **Icons**: Bootstrap Icons
- **Notifications**: Toastr

## 🚀 Getting Started

### Prerequisites

- Go 1.21+ installed
- PostgreSQL installed and running
- Web browser

### Backend Setup

1. **Navigate to backend directory**

   ```bash
   cd backend
   ```

2. **Create environment file**

   ```bash
   cp .env.example .env
   ```

3. **Configure `.env` file**

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=procurement_db
   JWT_SECRET=your-super-secret-jwt-key
   PORT=3000
   WEBHOOK_URL=https://webhook.site/your-unique-url
   ```

4. **Create database**

   ```sql
   CREATE DATABASE procurement_db;
   ```

5. **Install dependencies & run**

   ```bash
   go mod tidy
   go run main.go
   ```

   The API will be available at `http://localhost:3000`

### Frontend Setup

1. **Navigate to frontend directory**

   ```bash
   cd frontend
   ```

2. **Configure API URL** (optional)

   Edit `js/config.js` if your backend runs on a different port:

   ```javascript
   const API_BASE_URL = "http://localhost:3000/api";
   ```

3. **Run the frontend server**

   ```bash
   node server.js
   ```

   The frontend will be available at `http://localhost:8080`

4. **Open in browser**

   Navigate to `http://localhost:8080`

## 📚 API Documentation

### Authentication

| Method | Endpoint             | Description              |
| ------ | -------------------- | ------------------------ |
| POST   | `/api/auth/register` | Register new user        |
| POST   | `/api/auth/login`    | Login and get JWT token  |
| GET    | `/api/profile`       | Get current user profile |

### Items (Protected)

| Method | Endpoint         | Description     |
| ------ | ---------------- | --------------- |
| GET    | `/api/items`     | Get all items   |
| GET    | `/api/items/:id` | Get item by ID  |
| POST   | `/api/items`     | Create new item |
| PUT    | `/api/items/:id` | Update item     |
| DELETE | `/api/items/:id` | Delete item     |

### Suppliers (Protected)

| Method | Endpoint             | Description         |
| ------ | -------------------- | ------------------- |
| GET    | `/api/suppliers`     | Get all suppliers   |
| GET    | `/api/suppliers/:id` | Get supplier by ID  |
| POST   | `/api/suppliers`     | Create new supplier |
| PUT    | `/api/suppliers/:id` | Update supplier     |
| DELETE | `/api/suppliers/:id` | Delete supplier     |

### Purchases (Protected)

| Method | Endpoint             | Description         |
| ------ | -------------------- | ------------------- |
| GET    | `/api/purchases`     | Get all purchases   |
| GET    | `/api/purchases/:id` | Get purchase by ID  |
| POST   | `/api/purchases`     | Create new purchase |

### Request/Response Examples

**Login Request:**

```json
POST /api/auth/login
{
    "username": "admin",
    "password": "password123"
}
```

**Login Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "user"
    }
  }
}
```

**Create Purchase Request:**

```json
POST /api/purchases
Authorization: Bearer <token>
{
    "supplier_id": 1,
    "items": [
        {"item_id": 1, "qty": 5},
        {"item_id": 2, "qty": 3}
    ]
}
```

## ✨ Features

### Backend

- ✅ User authentication (Register & Login)
- ✅ JWT token-based authorization
- ✅ Password hashing with bcrypt
- ✅ CRUD operations for Items & Suppliers
- ✅ Purchase transaction with ACID compliance (database transaction)
- ✅ Server-side calculation of SubTotal & GrandTotal
- ✅ Stock validation and automatic deduction
- ✅ Webhook notification after successful purchase
- ✅ Input validation
- ✅ CORS enabled

### Frontend

- ✅ Login & Register pages
- ✅ JWT token handling (LocalStorage)
- ✅ Dashboard with statistics
- ✅ Items management (CRUD)
- ✅ Suppliers management (CRUD)
- ✅ Shopping cart functionality (client-side)
- ✅ Event delegation for dynamic elements
- ✅ Reusable AJAX wrapper with automatic auth header
- ✅ Toast notifications for user feedback
- ✅ Responsive design with Bootstrap 5

## 🔐 Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt
2. **JWT Authentication**: Secure token-based authentication
3. **Protected Routes**: Middleware to protect sensitive endpoints
4. **Input Validation**: Server-side validation for all inputs
5. **XSS Prevention**: HTML escaping on frontend

## 💡 Bonus Features Implemented

### Backend Bonus

- ✅ **Database Transaction (ACID)**: Purchase creation uses transaction with rollback on failure
- ✅ **Webhook Integration**: HTTP POST notification sent after successful purchase

### Frontend Bonus

- ✅ **Event Delegation**: Properly implemented for dynamically created buttons
- ✅ **Reusable AJAX**: Modular `api` object with automatic auth header injection
- ✅ **Robust Error Handling**: Toastr notifications for all error scenarios

## 📝 Database Schema

```
Users
├── ID (PK)
├── Username (Unique)
├── Password (Hashed)
├── Role
└── Timestamps

Suppliers
├── ID (PK)
├── Name
├── Email
├── Address
└── Timestamps

Items
├── ID (PK)
├── Name
├── Stock
├── Price
└── Timestamps

Purchasings
├── ID (PK)
├── Date
├── SupplierID (FK → Suppliers)
├── UserID (FK → Users)
├── GrandTotal
└── Timestamps

PurchasingDetails
├── ID (PK)
├── PurchasingID (FK → Purchasings)
├── ItemID (FK → Items)
├── Qty
├── SubTotal
└── Timestamps
```

## 🧪 Testing

### Quick Test Flow

1. **Register** a new user account
2. **Login** with the registered account
3. **Add Items**: Navigate to Items page and create some items
4. **Add Suppliers**: Navigate to Suppliers page and create suppliers
5. **Create Purchase**:
   - Select a supplier
   - Add items to cart
   - Submit the order
6. **View History**: Check purchase history

### API Testing with cURL

```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Get Items (with token)
curl http://localhost:3000/api/items \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 📄 License

This project is created for technical test purposes.

---

**Author**: Technical Test Submission  
**Date**: December 2025
