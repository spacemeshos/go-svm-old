package svm

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_RuntimeSuccess(t *testing.T) {
	rt, err := NewRuntime()
	assert.NoError(t, err)
	defer rt.Destroy()

	// DEPLOY

	code, err := ioutil.ReadFile("wasm/counter.wasm")
	assert.NoError(t, err)
	dplyMsg, err := DeployMessage{
		WasmFixedGasCode(code),
		Ctors{"initialize"},
		FixedLayout{4},
	}.Encode()
	assert.NoError(t, err)
	defer dplyMsg.Destroy()

	err = rt.ValidateDeploy(dplyMsg)
	assert.NoError(t, err)

	dplyEnv := DefaultEnvelope(StringAddress("@"))
	defer dplyEnv.Destroy()
	dplyCtx := NewContext()
	defer dplyCtx.Destroy()
	dplyRcpt, err := rt.Deploy(dplyEnv, dplyMsg, dplyCtx)
	assert.NoError(t, err)
	assert.True(t, dplyRcpt.Success)

	// SPAWN

	spwnMsg, err := SpawnMessage{
		dplyRcpt.Version,
		dplyRcpt.Address,
		"My Account",
		"initialize",
		Call32u(11).Encode(),
	}.Encode()
	assert.NoError(t, err)
	defer spwnMsg.Destroy()

	err = rt.ValidateSpawn(spwnMsg)
	assert.NoError(t, err)

	spwnCtx := NewContext()
	defer spwnCtx.Destroy()
	spwnEvn := DefaultEnvelope(StringAddress("@"))
	defer spwnEvn.Destroy()
	spwnRcpt, err := rt.Spawn(spwnEvn, spwnMsg, spwnCtx)
	assert.NoError(t, err)
	assert.True(t, spwnRcpt.Success)

	// CALL

	callMsg, err := CallMessage{
		dplyRcpt.Version,
		spwnRcpt.Address,
		"add",
		nil,
		Call32u(5).Encode(),
	}.Encode()
	assert.NoError(t, err)
	defer callMsg.Destroy()

	err = rt.ValidateCall(callMsg)
	assert.NoError(t, err)

	callCtx := spwnRcpt.Context()
	defer callCtx.Destroy()
	callEnv := DefaultEnvelope(StringAddress("@"))
	defer callEnv.Destroy()
	callRcpt, err := rt.Call(callEnv, callMsg, callCtx)
	assert.NoError(t, err)
	assert.True(t, callRcpt.Success)
	d, err := callRcpt.Return.Decode()
	assert.NoError(t, err)
	assert.Equal(t, 11, int(d[0].Uint()))
	assert.Equal(t, 16, int(d[1].Uint()))
}

func Test_RuntimeFailure(t *testing.T) {
	rt, err := NewRuntime()
	assert.NoError(t, err)
	defer rt.Destroy()

	// DEPLOY

	code, err := ioutil.ReadFile("wasm/failure.wasm")
	assert.NoError(t, err)
	dplyMsg, err := DeployMessage{
		WasmFixedGasCode(code),
		Ctors{"initialize"},
		FixedLayout{4},
	}.Encode()
	assert.NoError(t, err)
	defer dplyMsg.Destroy()

	err = rt.ValidateDeploy(dplyMsg)
	assert.NoError(t, err)

	dplyEnv := DefaultEnvelope(StringAddress("@"))
	defer dplyEnv.Destroy()
	dplyCtx := NewContext()
	defer dplyCtx.Destroy()
	dplyRcpt, err := rt.Deploy(dplyEnv, dplyMsg, dplyCtx)
	assert.NoError(t, err)
	assert.True(t, dplyRcpt.Success)

	// SPAWN

	spwnMsg, err := SpawnMessage{
		dplyRcpt.Version,
		dplyRcpt.Address,
		"My Account",
		"initialize",
		Call32u(11).Encode(),
	}.Encode()
	assert.NoError(t, err)
	defer spwnMsg.Destroy()

	err = rt.ValidateSpawn(spwnMsg)
	assert.NoError(t, err)

	spwnCtx := NewContext()
	defer spwnCtx.Destroy()
	spwnEvn := DefaultEnvelope(StringAddress("@"))
	defer spwnEvn.Destroy()
	spwnRcpt, err := rt.Spawn(spwnEvn, spwnMsg, spwnCtx)
	assert.NoError(t, err)
	assert.True(t, spwnRcpt.Success)

	// CALL

	callMsg, err := CallMessage{
		dplyRcpt.Version,
		spwnRcpt.Address,
		"fail",
		nil,
		nil,
	}.Encode()
	assert.NoError(t, err)
	defer callMsg.Destroy()

	err = rt.ValidateCall(callMsg)
	assert.NoError(t, err)

	callCtx := spwnRcpt.Context()
	defer callCtx.Destroy()
	callEnv := DefaultEnvelope(StringAddress("@"))
	defer callEnv.Destroy()
	callRcpt, err := rt.Call(callEnv, callMsg, callCtx)
	assert.NoError(t, err)
	assert.False(t, callRcpt.Success)
}
