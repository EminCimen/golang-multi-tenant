# Multi-Tenant API with Go

A multi-tenant REST API built with Go, featuring separate database isolation for each tenant. This project demonstrates a secure and scalable approach to building multi-tenant applications.

## Features

- 🔐 Multi-tenant architecture with database isolation
- 📚 Swagger documentation
- 🔑 JWT authentication
- 👥 User management per tenant
- 📝 Post management system
- 🛡️ Secure password hashing
- 🎯 Clean project structure

## Tech Stack

- [Go](https://golang.org/) - Programming language
- [Gin](https://gin-gonic.com/) - Web framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [JWT](https://github.com/golang-jwt/jwt) - Authentication
- [Swagger](https://swagger.io/) - API documentation

## Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Make sure PostgreSQL is running and accessible

## Installation

1. Clone the repository:

```bash
git clone https://github.com/EminCimen/golang-multi-tenant.git
cd golang-multi-tenant
```

2. Install dependencies:

```bash
go mod download
```

3. Create `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-key-here
JWT_EXPIRATION_HOURS=24
```

4. Run the application:

```bash
go run main.go
```

The server will start at `http://localhost:8080`

## API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

### Main Endpoints

- POST `/tenants` - Create a new tenant
- POST `/register` - Register a new user for a tenant
- POST `/login` - Login user
- GET `/me` - Get current user info
- POST `/posts` - Create a new post
- GET `/posts` - List all posts
- GET `/posts/{id}` - Get a specific post

## Project Structure

```
.
├── internal/
│   ├── api/         # API handlers
│   ├── database/    # Database configuration and connections
│   ├── middleware/  # Middleware functions
│   └── models/      # Data models
├── docs/           # Swagger documentation
├── main.go        # Application entry point
├── go.mod         # Go modules file
├── go.sum         # Go modules checksum
└── .env           # Environment variables
```

## Multi-Tenant Architecture

This project implements a database-per-tenant architecture:

1. Each tenant gets their own PostgreSQL database
2. Databases are named in the format: `tenant_[tenant_name]`
3. A main database (`tenant_management`) keeps track of all tenants
4. Complete data isolation between tenants
5. Shared authentication system with tenant-specific user management

## Authentication Flow

1. Create a tenant:

```json
POST /tenants
{
    "name": "Example Company"
}
```

2. Register a user:

```json
POST /register
{
    "tenant_id": 1,
    "email": "user@example.com",
    "password": "password123"
}
```

3. Login to get JWT token:

```json
POST /login
{
    "tenant_id": 1,
    "email": "user@example.com",
    "password": "password123"
}
```

4. Use the JWT token in the Authorization header:

```
Authorization: Bearer [your-jwt-token]
```

## Security

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Database-level tenant isolation
- Input validation and sanitization
- Secure password requirements

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Your Name - [@CimenDev](https://twitter.com/cimendev)
