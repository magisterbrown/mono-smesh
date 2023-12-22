import socket
import os

# Path to the Unix socket file
socket_path = './socket_test.sock'

# Create a Unix socket
server_socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

# Bind the socket to the path
server_socket.bind(socket_path)

# Listen for incoming connections
server_socket.listen(1)
print(f"Server is listening on {socket_path}")

while True:
    # Accept incoming connections
    connection, client_address = server_socket.accept()

    try:
        print('Connection established from', client_address)

        # Receive data from the client
        while True:
            data = connection.recv(1024)
            if not data:
                break
            print('Received:', data.decode())

            # Send back a response
            connection.sendall(b"Message received")

    finally:
        # Close the connection
        connection.close()
