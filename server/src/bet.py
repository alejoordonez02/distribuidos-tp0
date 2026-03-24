import datetime

# Simulated winner number in the lottery contest.
LOTTERY_WINNER_NUMBER = 7574


class Bet:
    """A lottery bet registry."""

    def __init__(
        self,
        agency: str,
        first_name: str,
        last_name: str,
        document: str,
        birthdate: str,
        number: str,
    ):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)


def has_won(bet: Bet) -> bool:
    """Checks whether a bet won the prize or not."""
    return bet.number == LOTTERY_WINNER_NUMBER


# old
class Person:
    def __init__(self, name: str, surname: str, birth: str):
        self.name = name
        self.surnname = surname
        self.birth = birth


class BetOld:
    def __init__(self, number: int, person: Person):
        self.number = number
        self.person = person
