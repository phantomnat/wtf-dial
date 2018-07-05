package internal_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/phantomnat/wtf-dial/internal/wtf"
	"github.com/phantomnat/wtf-dial/internal/wtf/bolt/internal"
)

// Ensure dial can be marshaled and unmarshaled.
func TestMarshalData(t *testing.T) {
	v := wtf.Dial{
		ID:      "ID",
		Token:   "TOKEN",
		Name:    "MyDial",
		Level:   10.2,
		ModTime: time.Now().UTC(),
	}

	var other wtf.Dial
	if buf, err := internal.MarshalDial(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalDial(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
