package curry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_ToString(t *testing.T) {
	// this is mostly a placeholder test to make sure everything is working as intended during repo setup
	t.Parallel()
	query := Base("test")

	assert.Equal(t, "test", query.ToString())
}
