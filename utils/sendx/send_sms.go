package sendx

type SMS struct {
	reciver []string
	Content string
}

func (s *SMS) AddReceiver(reciver ...string) Sender {
	s.reciver = append(s.reciver, reciver...)
	return s
}

func (s *SMS) AddContent(content string) Sender {
	s.Content = content
	return s
}

func (s *SMS) Send() error {
	return nil
}
