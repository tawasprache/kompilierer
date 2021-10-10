package langserv

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/alecthomas/repr"
	"github.com/sourcegraph/jsonrpc2"
)

type method func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{}
type methodMap map[string]method

func aus(i interface{}) {
	log.Println(repr.String(i))
}

func zu(fn interface{}) func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{} {
	val := reflect.ValueOf(fn)
	in := val.Type().In(2)
	return func(ctx context.Context, conn jsonrpc2.JSONRPC2, params json.RawMessage) interface{} {
		v := reflect.New(in)
		json.Unmarshal(params, v.Interface())
		ret := val.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(conn), v.Elem()})
		if len(ret) == 0 {
			return nil
		} else {
			if !ret[0].IsNil() {
				return ret[0].Interface()
			}
			if !ret[1].IsNil() {
				return ret[1].Interface()
			}
			panic("e")
		}
	}
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}

func StartServer() {
	s := server{}
	a := methodMap{
		"initialize":             zu(s.Initialize),
		"initialized":            zu(s.Initialized),
		"textDocument/didOpen":   zu(s.DidOpen),
		"textDocument/didChange": zu(s.DidChange),
		"textDocument/didClose":  zu(s.DidClose),
	}
	han := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
		v, ok := a[req.Method]
		if !ok {
			return nil, errors.New("not found")
		}
		resp := v(ctx, conn, *req.Params)

		return resp, nil
	})
	<-jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{}), han).DisconnectNotify()
}
