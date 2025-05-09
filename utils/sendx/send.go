package sendx

type Sender interface {
	AddReceiver(reciver ...string) Sender
	SetContent(content string) Sender
	Send() error
}
