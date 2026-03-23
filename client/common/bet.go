package common

const LEN_BIRTH = 10

type Person struct {
	Name    string
	Surname string
	birth   [LEN_BIRTH]byte
}

func (p *Person) Birth() string {
	s := string(p.birth[:])
	return s
}

type Bet struct {
	Number uint64
	Person Person
}

func newPerson(name string, surname string, birth string) Person {
	var _birth [LEN_BIRTH]byte
	copy(_birth[:], birth)

	person := Person{Name: name, Surname: surname, birth: _birth}
	return person
}

func NewBet(Number uint64, name string, surname string, birth string) Bet {
	Person := newPerson(name, surname, birth)

	bet := Bet{Number, Person}
	return bet
}
