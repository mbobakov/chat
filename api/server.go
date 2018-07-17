package api

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// server for the chat service
type server struct {
	clientsChans sync.Map // map[login]chan
	log          logger
}

type logger interface {
	Printf(string, ...interface{})
}

// Serve the chat
func Serve(ctx context.Context, listen string, options ...option) error {
	s := &server{
		log: log.New(ioutil.Discard, "", 0),
	}
	for _, o := range options {
		o(s)
	}

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return errors.Wrapf(err, "Couldn't start listen on the interface: '%s'", listen)
	}

	grpc_prometheus.EnableHandlingTimeHistogram(func(h *prometheus.HistogramOpts) {})

	mwu := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
	}
	mws := []grpc.StreamServerInterceptor{
		grpc_prometheus.StreamServerInterceptor,
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(mwu...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(mws...)),
	)
	RegisterChatServer(srv, s)
	s.log.Printf("[INFO] Start listening on '%s'", listen)
	e := make(chan error, 1)
	go func() {
		e <- srv.Serve(lis)
	}()
	select {
	case <-ctx.Done():
		s.log.Printf("[INFO] Server exiting... ")
		return nil
	case err := <-e:
		return err
	}
}

//go:generate charlatan -package charlatan -output subsribrServer_mock_test.go Chat_SubscribeServer
func (s *server) Subscribe(stream Chat_SubscribeServer) error {
	user, err := userFromMetadata(stream.Context())
	if err != nil {
		return errors.Wrap(err, "Couldn't get user from connection")
	}
	_, ok := s.clientsChans.Load(user)
	if ok {
		return errors.Errorf("User '%s' is already registred", user)
	}
	if user == "//" {
		return errors.Errorf("'%s' is internal server name", "//")
	}

	defer s.log.Printf("[INFO] Stream is end for %s", user)
	defer s.broadcast("//", []byte(user+" is diconnected"))
	defer s.clientsChans.Delete(user)
	gr, ctx := errgroup.WithContext(stream.Context())

	gr.Go(func() error {
		return s.listenFrom(ctx, user, stream)
	})

	toClientChan := make(chan *Message)
	gr.Go(func() error {
		return s.listenTo(ctx, stream, toClientChan)
	})
	s.broadcast("//", []byte(user+" is connected"))
	s.clientsChans.Store(user, toClientChan)
	s.log.Printf("[INFO] User '%s' is connected", user)
	return gr.Wait()
}

func (s *server) broadcast(from string, msg []byte) {
	s.clientsChans.Range(func(login, ch interface{}) bool {
		user := login.(string)
		channel := ch.(chan *Message)
		channel <- &Message{Body: msg, To: user, From: from}
		return true
	})
}

func userFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || md.Len() == 0 {
		return "", errors.New("No username in incoming metadata")
	}
	uu := md.Get(userContextKey)
	if len(uu) == 0 {
		return "", errors.New("No username in incoming metadata")
	}
	return uu[0], nil
}

func (s *server) listenFrom(ctx context.Context, f string, stream Chat_SubscribeServer) error {
	fromClient := make(chan *Message)
	e := make(chan error)

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				e <- nil
				return
			}
			if err != nil {
				e <- err
				return
			}
			in.From = f
			fromClient <- in
		}
	}()
	for {
		select {
		case m := <-fromClient:
			s.log.Printf("[DEBUG] '%s' send message '%s'", m.From, m.Body)
			s.broadcast(f, m.Body)
		case <-ctx.Done():
			return nil
		case err := <-e:
			return err
		}
	}
}

func (s *server) listenTo(ctx context.Context, stream Chat_SubscribeServer, toClient <-chan *Message) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-toClient:
			s.log.Printf("[DEBUG] Send a message '%s' from '%s' to '%s'", m.Body, m.From, m.To)
			err := stream.Send(m)
			if err != nil {
				return err
			}
		}
	}
}
