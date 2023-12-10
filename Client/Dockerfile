# Use an official Go runtime as a parent image
FROM golang:1.21.4

# Set the working directory in the container
WORKDIR /app

# Copy the entire project directory into the container
COPY . .


# Build the Go applications inside the container
RUN go build -o Chordexe

# Specify the commands to run your binaries
CMD ["./Chordexe", "-a", "127.0.0.1", "-p", "1234", "--ts", "30000", "--tff", "10000", "--tcp", "30000", "-r", "4", "-i", "1"]




