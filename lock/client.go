package lock

type LockClient interface {
	Init() error
	Get() (*Semaphore, error)
	Set(*Semaphore) error
}
