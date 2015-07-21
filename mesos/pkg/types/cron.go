package types

// Cron struct executes code with a time interval
type Cron struct {
	Frequency string
	Code      func()
}
