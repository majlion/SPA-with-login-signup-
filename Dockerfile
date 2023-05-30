# FROM golang:1.16-alpine

# WORKDIR /app

# COPY ./backend .

# RUN go build -o main .

# CMD ["./main"]



FROM golang:1.16-alpine

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy Go module files
COPY ./backend .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Set up SQLite database file
RUN mkdir -p db
RUN touch db/mydatabase.db

CMD ["./main"]
