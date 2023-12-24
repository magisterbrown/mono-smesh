import random

def agent(state: dict) -> dict:
    res = random.choice([0,1,2])
    #if state["observation"] == 1:
    #    raise ValueError
    return res

if __name__=='__main__':
    print(agent({}))
