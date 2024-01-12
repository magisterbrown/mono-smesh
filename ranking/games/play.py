import argparse
import socket
import json
from typing import Optional

import numpy as np
from fanorona_aec.env.fanorona_move import FanoronaMove
from fanorona_aec.env.fanorona_state import FanoronaState

# To test:
# nc -lkU ./test.sock
# python3 play.py ./test.sock

# Establish connection with a server that manages agents
parser = argparse.ArgumentParser()
parser.add_argument("socket")
args = parser.parse_args()

def request(data: dict, recv: int = 0) -> Optional[bytes]:
    with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as client:
        client.connect(args.socket)
        client.sendall(json.dumps(data).encode())
        if recv:
            return client.recv(recv)

def agn(nm):
    if nm is None:
        return ""
    return f"player_{int(nm)}"

# Game simulation
env = FanoronaState()
env.reset()
history = list()

while not env.done:
    cands = env.legal_moves
    board_state = str(env)
    history.append({"agent": agn(env.turn_to_play), "state": board_state})
    if len(cands)==1:
        env.push(FanoronaMove.from_action(cands[0]))
    else:
        action = request({"type": "move", "agent": agn(env.turn_to_play), "args": {"observation": board_state}}, 16)
        try:
            move = FanoronaMove.from_action(action)
            assert move.to_action() in env.legal_moves, "Illegal move"
        except:
            #TODO chack other
            request({"type": "done", "agent": agn(env.turn_to_play.other()), "broken": agn(env.turn_to_play)})
            raise AssertionError("TODO: handle illegal move")

        env.push(move)

request({"type": "done", "agent": agn(env.winner), "history": json.dumps(history)})
