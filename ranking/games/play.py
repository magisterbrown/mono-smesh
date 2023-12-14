from pettingzoo.classic import rps_v2
import json
import time
import getpass

env = rps_v2.env(max_cycles=1)
env.reset(seed=42)
acc_rewards = env.rewards.copy()

itt = env.agent_iter()
for agent in env.agent_iter():
    observation, reward, termination, truncation, info = env.last()
    if termination or truncation:
        action = None
    else:
        # this is where you would insert your policy
        print(json.dumps({"type": "move", "agent": agent, "observation": observation.tolist()}))
        decision = getpass.getpass(stream=open('/dev/null', 'w'))
        action = json.loads(decision)["choice"]

    env.step(action)
    for agent, reward in env.rewards.items():
        acc_rewards[agent]+=reward
env.close()
winner = max(acc_rewards, key=acc_rewards.get)
if all(value == 0 for value in acc_rewards.values()):
    winner = ""
print(json.dumps({"type": "result", "winner": winner}))
