package svm

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_RuntimeValidateDeploy(t *testing.T) {
	rt, err := NewRuntime()
	assert.NoError(t, err)
	defer rt.Destroy()

	// VALIDATE DEPLOY

	msgBytes, err := ioutil.ReadFile("test_assets/craft_deploy_example.bin")
	assert.NoError(t, err)
	msg := NewMessageFromBytes(msgBytes)
	defer msg.Destroy()

	err = rt.ValidateDeploy(msg)
	assert.NoError(t, err)
}
