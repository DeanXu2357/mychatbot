package probe

type Notifier interface {
	NotOK() <-chan struct{}
	OK() <-chan struct{}
}
