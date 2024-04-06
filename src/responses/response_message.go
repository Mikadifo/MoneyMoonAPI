package responses

type Message int

const (
	SUCCESS Message = iota
	PARTIAL Message = iota
	ERROR
)

func (message Message) String() string {
	return [...]string{"success", "partial", "error"}[message]
}
