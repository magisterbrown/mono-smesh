import argparse
import socket
import json

from pettingzoo.classic import rps_v2

# Establish connection with a server that manages agents
parser = argparse.ArgumentParser()
parser.add_argument("socket")
args = parser.parse_args()

client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
client.connect(args.socket)

# Game simulation

env = rps_v2.env(max_cycles=1)
env.reset(seed=42)
#acc_rewards = env.rewards.copy()
for agent in env.agent_iter():
    observation, reward, termination, truncation, info = env.last()
    if termination or truncation:
        action = None
    else:
        # this is where you would insert your policy
        client.sendall(json.dumps({"type": "move", "agent": agent, "args": {"observation": observation.tolist()}}).encode())
        action = client.recv(16)
    try:
        env.step(int(action))
    except:
        raise AssertionError("TODO: handle illegal move")
#    for agent, reward in env.rewards.items():
#        acc_rewards[agent]+=reward
env.close()

#winner = max(acc_rewards, key=acc_rewards.get)
#if all(value == 0 for value in acc_rewards.values()):
#    winner = ""
#print(json.dumps({"type": "result", "winner": winner}))
