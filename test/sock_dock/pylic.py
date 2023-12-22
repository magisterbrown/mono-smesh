import multiprocessing
import time
import socket 

pth = "./mysocket.sock"

def cl1():
    client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    client.connect(pth)
    msg = "cl1"
    client.sendall(msg.encode())
    time.sleep(2)
    print(f"cl1: {client.recv(4096)}")


def cl2():
    client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    client.connect(pth)
    time.sleep(0.5)
    msg = "cl2"
    #client.sendall(msg.encode())
    print(f"cl2: {client.recv(4096)}")

p1 = multiprocessing.Process(target=cl1)
p2 = multiprocessing.Process(target=cl2)

p1.start()
p2.start()
