import sys
import json
import random

def agent(state: dict) -> dict:
    res = random.choice([0,1,2])
    if state["observation"] == 1:
        raise ValueError
    return {"type":"decision", "choice": 0}

if __name__=='__main__':
    res = agent(json.loads(sys.argv[1]))
    print(json.dumps(res))
