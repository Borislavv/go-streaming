package enum

const (
	StreamByID           Actions = "ID"
	StreamByIDWithOffset Actions = "ID_WITH_OFFSET"
)

type Actions string

func (a Actions) String() string {
	return string(a)
}
