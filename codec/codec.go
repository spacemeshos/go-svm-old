package codec

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
	"go-svm/common"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type (
	ReceiptDeployTemplate = common.ReceiptDeployTemplate
	ReceiptSpawnApp       = common.ReceiptSpawnApp
	ReceiptExecApp        = common.ReceiptExecApp
)

const (
	markerErr = 0
	markerOk  = 1
)

var (
	instance *wasmer.Instance
)

func init() {
	var err error
	instance, err = newInstance(codecWasmFilePath())
	if err != nil {
		panic(err)
	}
}

func codecWasmFilePath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(file, "../../svm/svm_codec.wasm")
}

func newInstance(filename string) (*wasmer.Instance, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	module, err := wasmer.NewModule(store, bytes)
	if err != nil {
		return nil, err
	}

	instance, err := wasmer.NewInstance(module, wasmer.NewImportObject())
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func EncodeTxDeployTemplate(version int, name string, code []byte, data []byte) ([]byte, error) {
	txJson, err := json.Marshal(struct {
		Version int    `json:"version"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Data    string `json:"data"`
	}{
		Version: version,
		Name:    name,
		Code:    hex.EncodeToString(code),
		Data:    hex.EncodeToString(data),
	})
	if err != nil {
		return nil, err
	}

	argPtr, err := newBuffer(txJson)
	if err != nil {
		return nil, err
	}

	fn, err := instance.Exports.GetFunction("wasm_deploy_template")
	if err != nil {
		return nil, err
	}

	retPtr, err := fn(argPtr)
	if err != nil {
		return nil, err
	}

	return loadBuffer(retPtr.(int32))
}

func EncodeTxSpawnApp(version int, templateAddr []byte, name string, ctorName string, calldata []byte) ([]byte, error) {
	txJson, err := json.Marshal(struct {
		Version      int    `json:"version"`
		TemplateAddr string `json:"template"`
		Name         string `json:"name"`
		CtorName     string `json:"ctor_name"`
		Calldata     string `json:"calldata"`
	}{
		Version:      version,
		TemplateAddr: hex.EncodeToString(templateAddr),
		Name:         name,
		CtorName:     ctorName,
		Calldata:     hex.EncodeToString(calldata),
	})
	if err != nil {
		return nil, err
	}

	argPtr, err := newBuffer(txJson)
	if err != nil {
		return nil, err
	}

	fn, err := instance.Exports.GetFunction("wasm_encode_spawn_app")
	if err != nil {
		return nil, err
	}

	retPtr, err := fn(argPtr)
	if err != nil {
		return nil, err
	}

	return loadBuffer(retPtr.(int32))
}

func EncodeTxExecApp(version int, appAddr []byte, funcName string, calldata []byte) ([]byte, error) {
	txJson, err := json.Marshal(struct {
		Version  int    `json:"version"`
		AppAddr  string `json:"app"`
		FuncName string `json:"func_name"`
		Calldata string `json:"calldata"`
	}{
		Version:  version,
		AppAddr:  hex.EncodeToString(appAddr),
		FuncName: funcName,
		Calldata: hex.EncodeToString(calldata),
	})
	if err != nil {
		return nil, err
	}

	argPtr, err := newBuffer(txJson)
	if err != nil {
		return nil, err
	}

	fn, err := instance.Exports.GetFunction("wasm_encode_exec_app")
	if err != nil {
		return nil, err
	}

	retPtr, err := fn(argPtr)
	if err != nil {
		return nil, err
	}

	return loadBuffer(retPtr.(int32))
}

func EncodeCallData(abi []string, data []int) ([]byte, error) {
	calldataJson, err := json.Marshal(struct {
		ABI  []string `json:"abi"`
		Data []int    `json:"data"`
	}{
		ABI:  abi,
		Data: data,
	})
	if err != nil {
		return nil, err
	}

	argPtr, err := newBuffer(calldataJson)
	if err != nil {
		return nil, err
	}

	fn, err := instance.Exports.GetFunction("wasm_encode_calldata")
	if err != nil {
		return nil, err
	}

	retPtr, err := fn(argPtr)
	if err != nil {
		return nil, err
	}

	ret, err := loadBuffer(retPtr.(int32))
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	if err := json.Unmarshal(ret, &v); err != nil {
		return nil, err
	}

	calldata, ok := v["calldata"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid json format: %s", ret)
	}

	bytes, err := hex.DecodeString(calldata)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %s", calldata)
	}

	return bytes, nil
}

func DecodeReturndata(rawReturndata []byte) (string, error) {
	calldataJson, err := json.Marshal(struct {
		Calldata string `json:"calldata"`
	}{
		Calldata: hex.EncodeToString(rawReturndata),
	})

	if err != nil {
		return "", err
	}

	argPtr, err := newBuffer(calldataJson)
	if err != nil {
		return "", err
	}

	fn, err := instance.Exports.GetFunction("wasm_decode_calldata")
	if err != nil {
		return "", err
	}

	retPtr, err := fn(argPtr)
	if err != nil {
		return "", err
	}

	ret, err := loadBuffer(retPtr.(int32))
	if err != nil {
		return "", err
	}

	var v map[string]interface{}
	if err := json.Unmarshal(ret, &v); err != nil {
		return "", err
	}

	return string(ret), nil
}

func DecodeReceiptDeployTemplate(rawReceipt []byte) (*ReceiptDeployTemplate, error) {
	v, err := decodeReceipt(rawReceipt)
	if err != nil {
		return nil, err
	}

	return v.(*ReceiptDeployTemplate), nil
}

func DecodeReceiptSpawnApp(rawReceipt []byte) (*ReceiptSpawnApp, error) {
	v, err := decodeReceipt(rawReceipt)
	if err != nil {
		return nil, err
	}

	return v.(*ReceiptSpawnApp), nil
}

func DecodeReceiptExecApp(rawReceipt []byte) (*ReceiptExecApp, error) {
	v, err := decodeReceipt(rawReceipt)
	if err != nil {
		return nil, err
	}

	return v.(*ReceiptExecApp), nil
}

func decodeReceipt(rawReceipt []byte) (interface{}, error) {
	decodeReceiptJson, err := json.Marshal(struct {
		Data string `json:"data"`
	}{
		Data: hex.EncodeToString(rawReceipt),
	})
	if err != nil {
		return nil, err
	}

	bufPtr, err := newBuffer(decodeReceiptJson)
	if err != nil {
		return nil, err
	}

	fn, err := instance.Exports.GetFunction("wasm_decode_receipt")
	if err != nil {
		return nil, err
	}

	retBufPtr, err := fn(bufPtr)
	if err != nil {
		return nil, err
	}

	ret, err := loadBuffer(retBufPtr.(int32))
	if err != nil {
		return nil, err
	}

	return decodeReceiptJSON(ret)

}

func decodeReceiptJSON(jsonReceipt []byte) (interface{}, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(jsonReceipt, &v); err != nil {
		return nil, err
	}

	errType, ok := v["err_type"].(string)
	if ok {
		switch errType {
		case "oog":
			return nil, errors.New("oog")
		case "template-not-found":
			templateAddr := v["template_addr"].(string)
			return nil, fmt.Errorf("template not found; template address: %v", templateAddr)
		case "app-not-found":
			appAddr := v["app_addr"].(string)
			return nil, fmt.Errorf("template not found; app address: %v", appAddr)
		case "compilation-failed":
			templateAddr := v["template_addr"].(string)
			appAddr := v["app_addr"].(string)
			msg := v["message"].(string)
			return nil, fmt.Errorf("compilation failed; template address: %v, app address: %v, msg: %v",
				templateAddr, appAddr, msg)
		case "instantiation-failed":
			templateAddr := v["template_addr"].(string)
			appAddr := v["app_addr"].(string)
			msg := v["message"].(string)
			return nil, fmt.Errorf("instantiation failed; template address: %v, app address: %v, msg: %v",
				templateAddr, appAddr, msg)
		case "function-not-found":
			templateAddr := v["template_addr"].(string)
			appAddr := v["app_addr"].(string)
			fnc := v["func"].(string)
			return nil, fmt.Errorf("function not found; template address: %v, app address: %v, func: %v",
				templateAddr, appAddr, fnc)
		case "function-failed":
			templateAddr := v["template_addr"].(string)
			appAddr := v["app_addr"].(string)
			fnc := v["func"].(string)
			msg := v["message"].(string)
			return nil, fmt.Errorf("function failed; template address: %v, app address: %v, func: %v, msg: %v",
				templateAddr, appAddr, fnc, msg)
		default:
			panic(fmt.Sprintf("invalid error type: %v", errType))
		}
	} else {
		ty := v["type"].(string)
		switch ty {
		case "deploy-template":
			success := v["success"].(bool)
			gasUsed := v["gas_used"].(float64)
			addr := v["addr"].(string)

			return &common.ReceiptDeployTemplate{
				Success:      success,
				TemplateAddr: common.BytesToAddress(mustDecodeHexString(addr)),
				GasUsed:      uint64(gasUsed),
			}, nil

		case "spawn-app":
			success := v["success"].(bool)
			app := v["app"].(string)
			state := v["state"].(string)
			returndata := v["returndata"].(string)
			logs := v["logs"].([]interface{})
			gasUsed := v["gas_used"].(float64)

			strLogs := make([]string, len(logs))
			for i, log := range logs {
				log := log.(map[string]interface{})
				code := log["code"].(float64)
				msg := log["msg"].(string)
				strLogs[i] = fmt.Sprintf("(code: %v, msg: %v)", code, msg)
			}

			return &common.ReceiptSpawnApp{
				Success:    success,
				AppAddr:    common.BytesToAddress(mustDecodeHexString(app)),
				State:      mustDecodeHexString(state),
				Returndata: mustDecodeHexString(returndata),
				Logs:       strLogs,
				GasUsed:    uint64(gasUsed),
			}, nil

		case "exec-app":
			success := v["success"].(bool)
			newState := v["new_state"].(string)
			returndata := v["returndata"].(string)
			logs := v["logs"].([]interface{})
			gasUsed := v["gas_used"].(float64)

			strLogs := make([]string, len(logs))
			for i, log := range logs {
				log := log.(map[string]interface{})
				code := log["code"].(float64)
				msg := log["msg"].(string)
				strLogs[i] = fmt.Sprintf("(code: %v, msg: %v)", code, msg)
			}

			return &common.ReceiptExecApp{
				Success:    success,
				NewState:   mustDecodeHexString(newState),
				Returndata: mustDecodeHexString(returndata),
				Logs:       strLogs,
				GasUsed:    uint64(gasUsed),
			}, nil

		default:
			panic(fmt.Sprintf("invalid receipt type: %v", ty))
		}
	}
}

func newBuffer(data []byte) (int32, error) {
	length := int32(len(data))
	ptr, err := bufferAlloc(length)
	if err != nil {
		return 0, err
	}

	bufferLength, err := bufferLength(ptr)
	if err != nil {
		return 0, err
	}
	if length != bufferLength {
		panic(fmt.Sprintf("allocated buffer size isn't sufficient; allocated: %v, got: %v", length, bufferLength))
	}

	dataPtr, err := bufferDataPtr(ptr)
	if err != nil {
		return 0, err
	}

	mem, err := instance.Exports.GetMemory("memory")
	if err != nil {
		return 0, err
	}

	copy(mem.Data()[dataPtr:], data)

	return ptr, nil
}

func loadBuffer(ptr int32) ([]byte, error) {
	length, err := bufferLength(ptr)
	if err != nil {
		return nil, err
	}

	dataPtr, err := bufferDataPtr(ptr)
	if err != nil {
		return nil, err
	}

	mem, err := instance.Exports.GetMemory("memory")
	if err != nil {
		return nil, err
	}

	buf := mem.Data()[dataPtr : dataPtr+length]
	marker := buf[0]
	data := buf[1:]

	switch marker {
	case markerErr:
		return nil, errors.New(string(data))
	case markerOk:
		return data, nil
	default:
		panic("invalid marker")
	}
}

func bufferAlloc(size int32) (int32, error) {
	fn, err := instance.Exports.GetFunction("wasm_alloc")
	if err != nil {
		return 0, err
	}

	buf, err := fn(size)
	return buf.(int32), err
}

func bufferLength(buf int32) (int32, error) {
	fn, err := instance.Exports.GetFunction("wasm_buffer_length")
	if err != nil {
		return 0, err
	}

	bufLen, err := fn(buf)
	return bufLen.(int32), err
}

func bufferDataPtr(buf int32) (int32, error) {
	fn, err := instance.Exports.GetFunction("wasm_buffer_data")
	if err != nil {
		return 0, err
	}

	dataPtr, err := fn(buf)
	return dataPtr.(int32), err
}

func mustDecodeHexString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return b
}
