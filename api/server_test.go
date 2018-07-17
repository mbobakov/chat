package api

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc/metadata"
)

func Test_server_subscribe(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"happyPath", false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			s := &server{
				log: log.New(ioutil.Discard, "", 0),
			}
			var (
				recvN int
			)
			stream := &FakeChat_SubscribeServer{
				ContextHook: func() context.Context {
					return metadata.NewIncomingContext(
						ctx,
						metadata.New(map[string]string{userContextKey: ""}),
					)
				},
				RecvHook: func() (*Message, error) {
					for recvN != 0 {
						time.Sleep(time.Second)
					}
					recvN++
					return &Message{}, nil
				},
				SendHook: func(*Message) error { return nil },
			}
			err := s.Subscribe(stream)
			time.Sleep(time.Second)
			cancel()

			assert.Equal(t, tt.wantErr, err != nil)
			assert.True(t, stream.ContextCalledN(2))
			assert.True(t, stream.RecvCalledN(2))
			assert.True(t, stream.SendCalledN(1))

		})
	}
}
