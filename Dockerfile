# Use the official Node.js 14 image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Install Node.js
RUN apt update && apt install -y nodejs npm

# Install the Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Install additional dependencies
RUN go install github.com/air-verse/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/google/wire/cmd/wire@latest

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install the dependencies
RUN npm install

# Copy the rest of the application code to the working directory
COPY . .

# Expose the port on which the application will run
EXPOSE 8080

# Build the application
RUN npm run build

# Start the application
CMD [ "npm", "start" ]