package sender

type Sender[a any] interface {
	Send(a) error
}
