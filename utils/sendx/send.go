package sendx

type Sender interface {
	AddReceiver(reciver ...string) Sender
	AddContent(content string) Sender
	Send() error
}
