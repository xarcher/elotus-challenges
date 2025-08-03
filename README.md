# Elotus Challenges

## Project Structure

```
elotus-challenges/
├── backend/          # Go backend server
├── frontend/         # Static frontend with Nginx
├── dsa/             # Data Structures and Algorithms challenges
├── setup/           # Database initialization scripts
├── files/           # Uploaded files storage
└── docker-compose.yml
```

## Installation and Running the Project

### System Requirements
- Docker and Docker Compose
- Go 1.24+ (for local development)
- PostgreSQL 15+ (if not using Docker)

### 1. Clone Repository
```bash
git clone <repository-url>
cd elotus-challenges
```

### 2. Run with Docker Compose (Recommended)

#### Start all services:
```bash
docker-compose up -d
```

#### Check container status:
```bash
docker-compose ps
```

### 3. Access Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432

### 4. API Endpoints

#### Authentication
```bash
# Register new user
POST http://localhost:8080/api/auth/register
Content-Type: application/json
{
    "username": "testuser",
    "password": "password123"
}

# Login
POST http://localhost:8080/api/auth/login
Content-Type: application/json
{
    "username": "testuser",
    "password": "password123"
}
```

#### File Upload
```bash
# Upload file (requires authentication token)
POST http://localhost:8080/api/upload
Authorization: Bearer <your-jwt-token>
Content-Type: multipart/form-data
# Form data: file
```

## 🔧 Configuration

### Environment Variables
Main configuration in `backend/config/config.yml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "postgres"      # Change to localhost for local dev
  port: 5432
  user: "elotus"
  password: "elotus_password"
  dbname: "elotus_test"
  sslmode: "disable"

jwt:
  secret_key: "changeit"  # Change in production
  expires_in: "24h"

upload:
  max_file_size: 8388608  # 8MB
  temp_dir: "./files"
```

### Database Configuration
Database will be automatically initialized with schema from `setup/sql-init.sql`:
- Users table
- File uploads table  
- Revoked tokens table