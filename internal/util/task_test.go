package util_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/csh0101/netagent.git/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestTaskExecutor(t *testing.T) {

	e := util.GetTaskExecutor()

	var flag = false

	finished, err := e.RunTask(context.TODO(), "test_task", func() error {
		if !flag {
			flag = true
			return errors.New("trap error")
		}
		time.Sleep(time.Second * 5)
		return nil
	})

	assert.Nil(t, err)

	if err != nil {
		assert.Nil(t, finished)
	}

	<-finished

}

// todo add other test for executor
