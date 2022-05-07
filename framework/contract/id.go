package contract

const IDKey = "nade:id"

type IDService interface {
	NewID() string
}
