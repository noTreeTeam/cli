package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

func TestDeleteAll(t *testing.T) {
	service := "test-cli"
	// Nothing to delete
	err := keyring.DeleteAll(service)
	if err != nil && err.Error() == "The name org.freedesktop.secrets was not provided by any .service files" {
		t.Skip("Keyring service not available in this environment")
	}
	assert.NoError(t, err)
	// Setup 2 items
	err = keyring.Set(service, "key1", "value")
	assert.NoError(t, err)
	err = keyring.Set(service, "key2", "value")
	assert.NoError(t, err)
	// Delete all items
	err = keyring.DeleteAll(service)
	assert.NoError(t, err)
	// Check items are gone
	_, err = keyring.Get(service, "key1")
	assert.ErrorIs(t, err, keyring.ErrNotFound)
	_, err = keyring.Get(service, "key2")
	assert.ErrorIs(t, err, keyring.ErrNotFound)
}
