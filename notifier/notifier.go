package notifier

type Notifier interface {
	Notify(interface{})
}

type NotifyService struct {

}

func (notifier *NotifyService) Notify(interface{}) {

}