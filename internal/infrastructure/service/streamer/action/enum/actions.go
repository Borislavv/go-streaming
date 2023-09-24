package enum

const (
	StreamByID Actions = "ID"
)

type Actions string

func (a Actions) String() string {
	return string(a)
}
