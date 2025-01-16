package models

type ErrServiceFailure struct {
	ServiceName string
}

func (e ErrServiceFailure) Error() string {
	return e.ServiceName + " service is not available"
}

type ErrRateLimitExceeded struct {
	ServiceName string
}

func (e ErrRateLimitExceeded) Error() string {
	return e.ServiceName + " service rate limit exceeded"
}
