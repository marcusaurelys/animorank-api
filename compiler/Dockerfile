# Use a lightweight image with GCC installed
FROM gcc:latest

# Set the working directory
WORKDIR /app

# Copy the compiler program into the container
COPY compiler.go .

# Install Go to run the `compiler.go` file
RUN apt-get update && apt-get install -y golang && \
    rm -rf /var/lib/apt/lists/*

# Build the compiler.go into a binary
RUN go build -o compiler compiler.go

# Set the default command to run the compiled binary
CMD ["./compiler"]
