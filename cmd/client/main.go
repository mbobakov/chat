package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"

	flags "github.com/jessevdk/go-flags"
	"github.com/mbobakov/chat/api"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ErrCanceled is a error for gracefull client shutdown
var ErrCanceled = errors.New("Canceled")

func main() {
	var opts = struct {
		Server   string `long:"server" env:"CHAT_SERVER" default:"127.0.0.1:8080" description:"Chat server for connecting to"`
		Username string `long:"username" env:"CHAT_USERNAME" default:"johhDoe" description:"User"`
		Password string `long:"password" env:"CHAT_PASSWORD"  default:"12345678" description:"Password"`
	}{}
	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	for err != ErrCanceled {
		var (
			cc     *grpc.ClientConn
			stream api.Chat_SubscribeClient
		)
		color.New(color.FgYellow, color.Italic).Printf("Connecting to %s\n", opts.Server)
		cc, err = grpc.Dial(opts.Server, grpc.WithInsecure())
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		color.New(color.FgYellow, color.Italic).Printf("Connected\n")
		client := api.NewChatClient(cc)
		stream, err = client.Subscribe(
			metadata.NewOutgoingContext(
				context.Background(),
				metadata.New(map[string]string{"chatusername": opts.Username}),
			),
		)
		if err != nil {
			color.New(color.FgYellow, color.Italic).Printf("Couldn't subscribe. Err: '%v'\n", err)
			cc.Close()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		msgCh := make(chan *api.Message, 10)
		grp, ctx := errgroup.WithContext(context.Background())
		grp.Go(func() error {
			return print(ctx, os.Stdout, msgCh)
		})
		grp.Go(func() error {
			return watchStream(ctx, stream, msgCh)
		})
		grp.Go(func() error {
			return watchStdin(ctx, os.Stdin, opts.Username, stream)
		})
		grp.Go(func() error {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-ctx.Done():
				return nil
			case <-sigs:
				return ErrCanceled
			}
		})
		color.New(color.FgYellow, color.Italic).Printf("Happy Chatting. Your Username: '%s'\n", opts.Username)

		err = grp.Wait()
		if err != nil && err != ErrCanceled {
			color.New(color.FgRed, color.Italic).Printf("Error: '%v'\n", err)
			cc.Close()
			os.Exit(1)
		}
		cc.Close()
	}

}

//go:generate charlatan -package charlatan -output test/charlatan/subcribrClient.go api.Chat_SubscribeClient
//go:generate charlatan -package charlatan -output test/charlatan/reader.go io.Reader
func watchStdin(ctx context.Context, in io.Reader, from string, stream api.Chat_SubscribeClient) error {
	scanner := bufio.NewScanner(in)
	o := make(chan []byte)
	e := make(chan error)
	go func() {
		for scanner.Scan() {
			o <- scanner.Bytes()
		}
		if scanner.Err() != nil {
			e <- scanner.Err()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-o:
			err := stream.Send(&api.Message{From: from, Body: m})
			if err != nil {
				return err
			}
		case err := <-e:
			return err
		}
	}
}

func watchStream(ctx context.Context, s api.Chat_SubscribeClient, out chan<- *api.Message) error {
	o := make(chan *api.Message)
	e := make(chan error)
	go func() {
		for {
			in, err := s.Recv()
			if err != nil {
				e <- err
				return
			}
			o <- in
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-o:
			out <- m
		case err := <-e:
			return err
		}
	}
}

func print(ctx context.Context, to io.Writer, in <-chan *api.Message) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case m := <-in:
			var userFn, msgFn func(string, ...interface{}) string
			userFn = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
			msgFn = color.New(color.FgGreen, color.Bold).SprintfFunc()
			if m.From == "//" {
				userFn = color.New(color.FgBlue, color.Italic).SprintfFunc()
				msgFn = color.New(color.FgCyan, color.Italic).SprintfFunc()
			}
			fmt.Fprintf(to, "%s %s\n", userFn("%s: ", m.From), msgFn("%s", m.Body))
		}
	}
}
