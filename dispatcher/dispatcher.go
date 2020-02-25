package dispatcher

type Dispatcher interface {
	Post(payload string, jobType string) (string, error)
}
