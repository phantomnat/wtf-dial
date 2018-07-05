package bolt_test

import (
	"reflect"
	"testing"

	"github.com/phantomnat/wtf-dial/internal/wtf"
)

// Ensure dial can be created and retrieved.
func TestDialService_CreateDial(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	dial := wtf.Dial{
		ID:    "xxx",
		Token: "yyy",
		Name:  "My Dial",
		Level: 50,
	}

	// Create new dial.
	if err := s.CreateDial(&dial); err != nil {
		t.Fatal(err)
	}

	// Retrieve dial and compare.
	other, err := s.Dial("xxx")
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&dial, other) {
		t.Fatalf("unexpected dial: %#v", other)
	}
}

// Ensure dial validates the id.
func TestDialService_CreateDial_ErrDialIDRequired(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	if err := c.DialService().CreateDial(&wtf.Dial{ID: ""}); err != wtf.ErrDialIDRequired {
		t.Fatal(err)
	}
}

// Ensure duplicate dials are not allowed.
func TestDialService_CreateDial_ErrDialExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	if err := c.DialService().CreateDial(&wtf.Dial{ID: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := c.DialService().CreateDial(&wtf.Dial{ID: "x"}); err != wtf.ErrDialExist {
		t.Fatal(err)
	}
}

// Ensure dial's level can be updated
func TestDialService_SetLevel(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	// create new dials
	if err := s.CreateDial(&wtf.Dial{ID: "xxx", Token: "yyy", Level: 50}); err != nil {
		t.Fatal(err)
	} else if err := s.CreateDial(&wtf.Dial{ID: "aaa", Token: "bbb", Level: 80}); err != nil {
		t.Fatal(err)
	}

	// verify dial 1 created.
	if d, err := s.Dial("xxx"); err != nil {
		t.Fatal(err)
	} else if d.Level != 50 {
		t.Fatalf("unexpected dial #1 level: %f", d.Level)
	}

	// Update dial levels.
	if err := s.SetLevel("xxx", "yyy", 60); err != nil {
		t.Fatal(err)
	} else if err := s.SetLevel("aaa", "bbb", 10); err != nil {
		t.Fatal(err)
	}

	// verify dial 1 updated.
	if d, err := s.Dial("xxx"); err != nil {
		t.Fatal(err)
	} else if d.Level != 60 {
		t.Fatalf("unexpected dial #1 level: %f", d.Level)
	}

	// verify dial 2 updated.
	if d, err := s.Dial("aaa"); err != nil {
		t.Fatal(err)
	} else if d.Level != 10 {
		t.Fatalf("unexpected dial #2 level: %f", d.Level)
	}
}

// Ensure dial level cannot be updated if token doesn't match.
func TestDialService_ErrUnauthorized(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DialService()

	// Create new dial
	if err := s.CreateDial(&wtf.Dial{ID: "xxx", Token: "yyy", Level: 50}); err != nil {
		t.Fatal(err)
	}

	// Update dial level with wrong token.
	if err := s.SetLevel("xxx", "bad_token", 60); err != wtf.ErrUnauthorized {
		t.Fatal(err)
	}
}

// Ensure error is returned if dial doesn't exist.
func TestDialService_SetLevel_ErrDialNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	if err := c.DialService().SetLevel("xxx", "", 50); err != wtf.ErrDialNotFound {
		t.Fatal(err)
	}
}
