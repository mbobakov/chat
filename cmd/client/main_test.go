package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/mbobakov/chat/api"
	"github.com/mbobakov/chat/cmd/client/test/charlatan"
	"github.com/stretchr/testify/assert"
)

func Test_watchStdin(t *testing.T) {
	tests := []struct {
		name             string
		from             string
		input            []byte
		sendErr, scanErr error
		expected         *api.Message
		shouldErr        bool
	}{
		// TODO: Add test cases.
		{"happyPath", "f", []byte("input"), nil, nil, &api.Message{From: "f", Body: []byte("input")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.input)
			r := &charlatan.FakeReader{
				ReadHook: func(to []byte) (int, error) {
					if tt.scanErr != nil {
						return 0, tt.scanErr
					}
					return b.Read(to)
				},
			}
			var actualMsg *api.Message
			stream := &charlatan.FakeChat_SubscribeClient{
				SendHook: func(m *api.Message) error {
					actualMsg = m
					return tt.sendErr
				},
			}
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			go func() {
				err = watchStdin(ctx, r, tt.from, stream)
			}()
			time.Sleep(10 * time.Millisecond)
			assert.Equal(t, tt.expected, actualMsg)
			assert.Equal(t, tt.shouldErr, err != nil)
			assert.True(t, stream.SendCalledOnce())
			assert.True(t, r.ReadCalledN(2))
		})
	}
}

func Test_watchStream(t *testing.T) {
	tests := []struct {
		name     string
		inputMsg *api.Message
		inputErr error
		expected *api.Message
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"happyPath", &api.Message{From: "hp"}, nil, &api.Message{From: "hp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var recvN int
			stream := &charlatan.FakeChat_SubscribeClient{
				RecvHook: func() (*api.Message, error) {
					for recvN != 0 {
						time.Sleep(time.Second)
					}
					return tt.inputMsg, tt.inputErr
				},
			}
			var (
				err    error
				actual *api.Message
			)

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			out := make(chan *api.Message, 1)
			go func() {
				err = watchStream(ctx, stream, out)
			}()
			time.Sleep(10 * time.Millisecond)
			select {
			case actual = <-out:
			default:
				t.Error("No message reseived")
				return
			}
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.True(t, stream.RecvCalledN(2))
		})
	}
}
