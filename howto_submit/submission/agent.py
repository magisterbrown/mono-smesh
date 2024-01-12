import random
import argparse
from fanorona_aec.env.fanorona_move import FanoronaMove
from fanorona_aec.env.fanorona_state import FanoronaState

parser.add_argument("observation")
args = parser.parse_args()

def agent(state):
    res = random.choice(state.legal_moves)
    return res

if __name__=='__main__':
    print(agent(FanoronaState.set_from_board_str(args.observation)))
