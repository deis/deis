package types

type Cron struct {
	Frequency string
	Code      func()
}
