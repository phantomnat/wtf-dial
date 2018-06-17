package bolt

import (
	"github.com/phantomnat/wtf-dial/internal/wtf"
	"github.com/phantomnat/wtf-dial/internal/wtf/bolt/internal"
)

// Ensure DialService implements wtf.DialService
var _ wtf.DialService = &DialService{}

// DialService represents a serive for managing dials.
type DialService struct {
	session *Session
}

// Dial returns a dial by ID.
func (s *DialService) Dial(id wtf.DialID) (*wtf.Dial, error) {
	// Start read-only transaction.
	tx, err := s.session.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal dial.
	var d wtf.Dial
	if v := tx.Bucket([]byte("Dials")).Get(itob(int(id))); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalDial(v, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// CreateDial creates a new dial.
func (s *DialService) CreateDial(d *wtf.Dial) error {
	// Retrieve current session user.
	u, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	// Start read-write transaction.
	tx, err := s.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create new ID.
	b := tx.Bucket([]byte("Dials"))
	seq, _ := b.NextSequence()
	d.ID = wtf.DialID(seq)

	// Assign Dial to current user.
	d.UserID = u.ID
	d.ModTime = s.session.now

	// Marshal and insert record.
	if v, err := internal.MarshalDial(d); err != nil {
		return err
	} else if err := b.Put(itob(int(d.ID)), v); err != nil {
		return err
	}

	return tx.Commit()
}

// SetLevel sets the current WTF level for a dial.
func (s *DialService) SetLevel(id wtf.DialID, level float64) error {
	// Retrieve current session user.
	u, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	// Start read-write transaction.
	tx, err := s.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("Dials"))

	// Find and unmarshal dial.
	var d wtf.Dial
	if v := b.Get(itob(int(id))); v == nil {
		return wtf.ErrDialNotFound
	} else if err := internal.UnmarshalDial(v, &d); err != nil {
		return err
	}

	// Only allow owner to update level
	if d.UserID != u.ID {
		return wtf.ErrUnauthorized
	}

	// Update dial level.
	d.Level = level
	d.ModTime = s.session.now

	// Marshal and insert record.
	if v, err := internal.MarshalDial(&d); err != nil {
		return err
	} else if err := b.Put(itob(int(id)), v); err != nil {
		return err
	}

	return tx.Commit()
}
