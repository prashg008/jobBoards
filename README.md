# Job Boards API

The Job Boards API is a RESTful web service built with Go and the Gin framework. It provides endpoints for managing job postings, user authentication, and job applications.

## Assumptions

This project makes the following assumptions:

1. **Authentication**: Assumes users can register, login, and retrieve their own details. Uses JWT Bearer token for authentication.
2. **Job Posting**: Assumes employers can create, list, and delete job postings. Also, users can list and apply for jobs and view their applications.
3. **Database**: Assumes PostgreSQL database is used for storing job postings, users, and job applications.

## Setup

To set up and run the project locally, follow these steps:

1. **Clone the Repository:**
   ```bash
   git clone <repository_url>
   cd jobBoards
   ```

2. **Set Up Configuration:**
   - Create configuration files for different environments (`development.yaml`, `testing.yaml`, `production.yaml`) in the `config` directory.
   - Define database settings, JWT secret, and other configurations in each file.

3. **Database Setup:**
   - Install PostgreSQL and create databases according to the configurations.

4. **Run the Application:**
   ```bash
   go run main.go
   ```

5. **Access Swagger Documentation:**
   - Navigate to `http://localhost:8080/swagger/index.html` to view and interact with the API documentation.

## Environment Variables

The following environment variables can be set to configure the application:

- `ENVIRONMENT`: Specifies the environment (`development`, `testing`, `production`).
- Other environment-specific configuration variables can be defined in the configuration files.
