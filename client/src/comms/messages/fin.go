package messages

type Fin struct{}

func NewFin() Fin {
	fin := Fin{}
	return fin
}
