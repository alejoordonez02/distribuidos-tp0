from .net import Bet

# Simulated winner number in the lottery contest.
LOTTERY_WINNER_NUMBER = 7574


def has_won(bet: Bet) -> bool:
    """Checks whether a bet won the prize or not."""
    return bet.number == LOTTERY_WINNER_NUMBER
