package messages

type Bet struct {
	Agency    string
	FirstName string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

type BetBatch []Bet

func NewBet(
	Agency string,
	FirstName string,
	LastName string,
	Document string,
	BirthDate string,
	Number string) Bet {
	bet := Bet{Agency, FirstName, LastName, Document, BirthDate, Number}
	return bet
}

func (b *Bet) GetSize() int {
	betSize := len(b.FirstName) + len(b.LastName) +
		len(b.Document) + len(b.BirthDate) + len(b.Number)

	return betSize
}
