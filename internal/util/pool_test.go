package util_test

import (
	"testing"

	"github.com/csh0101/netagent.git/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestPoolGet(t *testing.T) {

	buf := util.GetBuf()

	assert.Equal(t, 1024, len(buf))

	util.PutBuf(buf)

}
