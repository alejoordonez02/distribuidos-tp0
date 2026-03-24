package common

type Bet struct {
	FirstName string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

type BetBatch []Bet

func NewBet(
	FirstName string,
	LastName string,
	Document string,
	BirthDate string,
	Number string) Bet {
	bet := Bet{FirstName, LastName, Document, BirthDate, Number}
	return bet
}
