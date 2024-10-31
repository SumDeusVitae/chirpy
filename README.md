# Chirpy Documentation

## Overview
Chirpy is a lightweight web application that allows users to create and manage "chirps" (short messages). It includes user authentication, metrics tracking, and file serving functionalities.

## Features
- User registration and login
- Create, read, update, and delete chirps
- Health checks and metrics endpoints
- File serving with cache control

## Dependencies
- Go (version 1.16+)
- PostgreSQL
- Required Go packages:
  - `github.com/{user_name}/chirpy/internal/database`
  - `github.com/google/uuid`
  - `github.com/joho/godotenv`
  - `github.com/lib/pq`

## Environment Variables
Chirpy requires the following environment variables to be set:

- `DB_URL`: Database connection string for PostgreSQL.
- `PLATFORM`: Application platform identifier.
- `SECRET`: Secret key for token signing.
- `POLKA_KEY`: API key for the Polka service.

## Setup Instructions

1. **Clone the repository**:
   ```bash
   git clone https://github.com/SumDeusVitae/chirpy.git
   cd chirpy
   ```

2. **Set up your PostgreSQL database**: Make sure to create a database and note down the connection string.

3. **Create a `.env` file** in the root directory with the necessary environment variables:
   ```plaintext
   DB_URL=your_database_url
   PLATFORM=your_platform
   SECRET=your_secret
   POLKA_KEY=your_polka_key
   ```

4. **Install dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the application**:
   ```bash
   go run main.go
   ```

   The application will start on `http://localhost:8080`.

## API Endpoints

### User Management
- **POST `/api/users`**: Create a new user.
- **POST `/api/login`**: Log in an existing user and receive tokens.
- **POST `/api/refresh`**: Refresh the authentication token.
- **PUT `/api/users`**: Update user information.
- **DELETE `/api/chirps/{chirpID}`**: Revoke a user's access.

### Chirps Management
- **POST `/api/chirps`**: Create a new chirp.
- **GET `/api/chirps`**: Retrieve all chirps.
- **GET `/api/chirps/{chirpID}`**: Retrieve a specific chirp by ID.
- **DELETE `/api/chirps/{chirpID}`**: Delete a chirp.

### Administration
- **GET `/admin/metrics`**: Retrieve application metrics.
- **POST `/admin/reset`**: Reset users or database state (admin only).

### Health Check
- **GET `/api/healthz`**: Check the health of the application.

### Webhooks
- **POST `/api/polka/webhooks`**: Handle incoming webhooks from Polka service.

## Middleware
Chirpy includes middleware for:
- **Cache Control**: Prevents browsers from caching responses.
- **Metrics Tracking**: Increments file server hits for monitoring.

## Data Models

### User
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "string",
  "token": "string",
  "refresh_token": "string",
  "is_chirpy_red": "boolean"
}
```

### Chirp
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "string",
  "user_id": "uuid"
}
```

## Conclusion
Chirpy is designed to be a simple yet effective platform for sharing quick thoughts. By following this documentation, you should be able to set up and run the application, as well as understand its functionalities and API.