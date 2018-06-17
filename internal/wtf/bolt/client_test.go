package bolt_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/phantomnat/wtf-dial/internal/wtf/bolt"
	"github.com/phantomnat/wtf-dial/internal/wtf/mock"
)

// Now is the mocked current time for testing.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Client is a test wrapper for bolt.Client.
type Client struct {
	*bolt.Client

	Authenticator mock.Authenticator
}

// NewClient returns a new instance of Client pointing at a temporary file.
func NewClient() *Client {
	// Generate temporary filename.
	f, err := ioutil.TempFile("", "wtf-bolt-client-")
	if err != nil {
		panic(err)
	}
	// defer f.Close()
	f.Close()

	// Create client wrapper.
	c := &Client{
		Client:        bolt.NewClient(),
		Authenticator: mock.DefaultAuthenticator(),
	}
	c.Path = f.Name()
	c.Now = func() time.Time { return Now }

	// Assign mocks to the implementation.
	c.Client.Authenticator = &c.Authenticator

	return c
}

// MustOpenClient returns a new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	defer os.Remove(c.Path)
	return c.Client.Close()
}
