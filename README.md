# TMR Backend

A Go-based backend service for TMR (Targeted Memory Reactivation) lab experiment management.

## Overview

This project provides a REST API service for managing TMR lab experiments, including subject management, lab test administration, and result tracking. It's built with Go and uses MySQL as the database.

## Tech Stack

- **Language**: Go 1.22.1
- **Framework**: Gin Web Framework
- **Database**: MySQL with GORM
- **Container**: Docker
- **CI/CD**: GitHub Actions
- **Cloud Services**: AWS (ECR, CodeDeploy, S3, EC2)

## Key Features

### 1. Subject Management

- Create new test subjects
- Verify subject existence
- Generate unique login IDs for subjects

### 2. Lab Test Management

- Support for pre-tests and post-tests
- Track breathing patterns during tests
- Record cue word presentations
- Store test results and histories

### 3. File Management

- Temporary file storage for test results
- Automatic file cleanup after expiration (60 minutes)
- Secure file download functionality

### 4. Slack Integration

- Automated notifications for test events
- Test result summaries with downloadable CSV files
- Real-time updates on experiment progress

## Project Structure

- config/: Configuration setup (DB, env)
- dto/: Data Transfer Objects
- entity/: Database entities
- handler/: HTTP request handlers
- model/: Business logic and database operations
- router/: Route definitions
- util/: Utility functions
- main.go: Application entry point

## API Endpoints

### Subject APIs

- POST /api/subjects - Create new subject
- GET /api/subjects/check - Check subject existence

### Lab APIs

- POST /api/labs/breathing - Record breathing history
- POST /api/labs/cue - Record cue presentation
- POST /api/labs/start-test - Start a new test session
- POST /api/labs/test - Record test history
- GET /api/labs/cue - Get target words

### File APIs

- GET /api/files/:filename - Download result files

### Health Check

- GET /api/health - Service health check

## Deployment

The project uses GitHub Actions for CI/CD pipeline:

1. Builds Docker image
2. Pushes to Amazon ECR
3. Deploys to EC2 using AWS CodeDeploy

### Prerequisites

- AWS credentials
- MySQL database
- Slack webhook URL
- Environment variables setup

### Environment Variables

- DB_HOST
- DB_USERNAME
- DB_PASSWORD
- DB_NAME
- SLACK_WEBHOOK_URL
- BASE_URL

## Development Setup

1. Clone the repository
2. Install dependencies: go mod download
3. Set up environment variables: Copy .env.example to .env and configure
4. Run the application: go run main.go

## Docker Support

Build the image:
docker build -t tmr-backend .

Run the container:
docker run -p 8080:8080 -v /path/to/temp_files:/temp_files tmr-backend
