package notifier

type Notifier interface {
	SendMessage(a string)
}

func Init() []Notifier {
	return []Notifier{
		newTelegram(),
	}
}
