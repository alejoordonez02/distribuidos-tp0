package comms

type Serializable interface {
	Serialize() []byte
}
