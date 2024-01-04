import argparse
import socket
import json
from typing import Optional

from pettingzoo.classic import rps_v2

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


# Game simulation
env = rps_v2.env(max_cycles=3)
env.reset(seed=42)
acc_rewards = env.rewards.copy()
history = list()
for agent in env.agent_iter():
    observation, reward, termination, truncation, info = env.last()
    acc_rewards[agent]+=reward
    if termination or truncation:
        break
    else:
        # this is where you would insert your policy
        action = request({"type": "move", "agent": agent, "args": {"observation": observation.tolist()}}, 16)
    try:
        actt = int(action)
        env.step(actt)
        history.append({"agent": agent, "move": actt})
    except:
        # TODO: handle move with error
        env.agents.remove(agent)
        request({"type": "done", "agent": env.agents[0], "broken": agent})
        raise AssertionError("TODO: handle illegal move")
env.close()

winner = max(acc_rewards, key=acc_rewards.get)
if all(value == 0 for value in acc_rewards.values()):
    winner = ""
request({"type": "done", "agent": winner, "history": json.dumps(history)})
