package service

// FactorialService handles factorial calculations
type FactorialService interface {
	ValidateNumber(number string) (int64, error)
}

type factorialService struct{}

// NewFactorialService creates a new factorial service
func NewFactorialService() FactorialService {
	return &factorialService{}
}

// ValidateNumber validates and parses the input number string
func (s *factorialService) ValidateNumber(number string) (int64, error) {
	panic("not implemented")
}
