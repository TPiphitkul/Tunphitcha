package policy

type Decision struct {
	RateLimit int  // requests per second (example)
	StepUpMFA bool
	Block     bool
}

func Decide(level string) Decision {
	switch level {
	case "high":
		return Decision{RateLimit: 1, StepUpMFA: true, Block: false}
	case "medium":
		return Decision{RateLimit: 10, StepUpMFA: false, Block: false}
	default:
		return Decision{RateLimit: 100, StepUpMFA: false, Block: false}
	}
}
