package mock

import (
	"github.com/phantomnat/wtf-dial/internal/app/wtf"
)

type DialService struct {
	DialFn      func(id wtf.DialID) (*wtf.Dial, error)
	DialInvoked bool

	CreateDialFn      func(dial *wtf.Dial) error
	CreateDialInvoked bool

	SetLevelFn      func(id wtf.DialID, level float64) error
	SetLevelInvoked bool
}

func (s *DialService) Dial(id wtf.DialID) (*wtf.Dial, error) {
	s.DialInvoked = true
	return s.DialFn(id)
}

func (s *DialService) CreateDial(dial *wtf.Dial) error {
	s.CreateDialInvoked = true
	return s.CreateDialFn(dial)
}

func (s *DialService) SetLevel(id wtf.DialID, level float64) error {
	s.SetLevelInvoked = true
	return s.SetLevelFn(id, level)
}

// UserService represents a service for managing user authentication.
type UserService struct {
	AuthenticateFn      func(token string) (*wtf.User, error)
	AuthenticateInvoked bool
}

func (s *UserService) Authenticate(token string) (*wtf.User, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(token)
}
