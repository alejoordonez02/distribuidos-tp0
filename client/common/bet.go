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
