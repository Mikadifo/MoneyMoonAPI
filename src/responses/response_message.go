package responses

type Message int

const (
	SUCCESS Message = iota
	ERROR
)

func (message Message) String() string {
	return [...]string{"success", "error"}[message]
}
