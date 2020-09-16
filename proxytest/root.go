package proxytest

import (
	"log"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// TODO: simulate OnQueueReady, OnTick

type RootFilterHost struct {
	baseHost
	context proxywasm.RootContext

	pluginConfiguration, vmConfiguration []byte
}

func NewRootFilterHost(ctx proxywasm.RootContext, pluginConfiguration, vmConfiguration []byte,
) (*RootFilterHost, func()) {
	host := &RootFilterHost{context: ctx,
		pluginConfiguration: pluginConfiguration,
		vmConfiguration:     vmConfiguration,
	}
	hostMux.Lock() // acquire the lock of host emulation
	rawhostcall.RegisterMockWASMHost(host)
	return host, func() {
		hostMux.Unlock()
	}
}

func (n *RootFilterHost) ConfigurePlugin() {
	size := len(n.pluginConfiguration)
	n.context.OnConfigure(size)
}

func (n *RootFilterHost) StartVM() {
	size := len(n.vmConfiguration)
	n.context.OnVMStart(size)
}

func (n *RootFilterHost) ProxyGetBufferBytes(bt types.BufferType, start int, maxSize int,
	returnBufferData **byte, returnBufferSize *int) types.Status {
	var buf []byte
	switch bt {
	case types.BufferTypePluginConfiguration:
		buf = n.pluginConfiguration
	case types.BufferTypeVMConfiguration:
		buf = n.vmConfiguration
	default:
		// delegate to baseHost
		return n.getBuffer(bt, start, maxSize, returnBufferData, returnBufferSize)
	}

	if start >= len(buf) {
		log.Printf("start index out of range: %d (start) >= %d ", start, len(buf))
		return types.StatusBadArgument
	}

	*returnBufferData = &buf[start]
	if maxSize > len(buf)-start {
		*returnBufferSize = len(buf) - start
	} else {
		*returnBufferSize = maxSize
	}
	return types.StatusOK
}