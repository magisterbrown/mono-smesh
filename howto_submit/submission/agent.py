import random
import argparse
from fanorona_aec.env.fanorona_move import FanoronaMove
from fanorona_aec.env.fanorona_state import FanoronaState

parser = argparse.ArgumentParser()
parser.add_argument("--observation")
args = parser.parse_args()

def agent(state):
    res = random.choice(state.legal_moves)
    return res

if __name__=='__main__':
    board = args.observation
    state = FanoronaState()
    print(agent(state.set_from_board_str(board)))
