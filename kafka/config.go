package kafka

type Config struct {
	Brokers        []string
	GroupID        string
	Topics         []string
	MinBytes       int
	MaxBytes       int
	CommitInterval int // ms
}
