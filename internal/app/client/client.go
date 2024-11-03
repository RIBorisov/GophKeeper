package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/RIBorisov/GophKeeper/pkg/server"
)

func say(a ...any) {
	fmt.Println(a...)
}

func sayf(format string, a ...any) {
	fmt.Printf(format, a...)
}

type Client struct {
	grpcClient pb.GophKeeperServiceClient
	localCache *sync.Map
	token      string
}

func NewClient(ctx context.Context) (*Client, error) {
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	return &Client{
		grpcClient: pb.NewGophKeeperServiceClient(conn),
		localCache: &sync.Map{},
		token:      "",
	}, nil
}

func (c *Client) Register(ctx context.Context, s *bufio.Scanner) (*pb.RegisterResponse, error) {
	say("Enter login")
	s.Scan()
	login := s.Text()
	say("Enter password")
	s.Scan()
	password := s.Text()
	in := &pb.RegisterRequest{Login: login, Password: password}
	resp, err := c.grpcClient.Register(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	c.token = resp.GetToken()

	return resp, nil
}

func (c *Client) LogIn(ctx context.Context, s *bufio.Scanner) (*pb.AuthResponse, error) {
	say("Enter login")
	s.Scan()
	login := s.Text()
	say("Enter password")
	s.Scan()
	password := s.Text()
	in := &pb.AuthRequest{Login: login, Password: password}
	resp, err := c.grpcClient.Auth(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to log in: %w", err)
	}

	c.token = resp.GetToken()

	return resp, nil
}

func (c *Client) GetData(ctx context.Context, s *bufio.Scanner) (*pb.GetResponse, error) {
	if c.token == "" {
		say("Unauthenticated user")
		return nil, ErrUserUnauthenticated
	}
	say("Input ID")
	s.Scan()
	input := s.Text()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("token", c.token))
	resp, err := c.grpcClient.Get(ctx, &pb.GetRequest{ID: input})
	if err != nil {
		return nil, err
	}

	if resp != nil {
		return resp, nil
	}

	sayf("Data ID: %s\n", resp.GetID())

	return nil, nil

}

func (c *Client) SaveData(ctx context.Context, s *bufio.Scanner) error {
	if c.token == "" {
		say("Unauthenticated user")
		return ErrUserUnauthenticated
	}
	in := buildInput(s)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("token", c.token))
	resp, err := c.grpcClient.Save(ctx, in)
	if err != nil {
		return err
	}

	if resp != nil {
		sayf("Saved ID: %s\n", resp.GetID())
	}

	return nil
}

func buildInput(s *bufio.Scanner) *pb.SaveRequest {
	say(ChooseKind)
	s.Scan()
	input := s.Text()
	in := &pb.SaveRequest{}
	switch strings.TrimSpace(input) {
	case "1":
		say("Number")
		s.Scan()
		number := s.Text()
		say("Month/Year")
		s.Scan()
		mmYY := s.Text()
		say("CVC/CVV")
		s.Scan()
		cvc := s.Text()

		in = &pb.SaveRequest{
			Kind: pb.Kind_CARD,
			Data: &pb.SaveRequest_Card{Card: &pb.Card{Number: number, MmYy: mmYY, Cvc: cvc}},
		}
	case "2":
		say("Enter text")
		s.Scan()
		text := s.Text()
		in = &pb.SaveRequest{
			Kind: pb.Kind_TEXT,
			Data: &pb.SaveRequest_Text{Text: text},
		}
	case "3":
		say("Enter login")
		s.Scan()
		login := s.Text()
		say("Enter password")
		s.Scan()
		password := s.Text()
		in = &pb.SaveRequest{
			Kind: pb.Kind_CREDENTIALS,
			Data: &pb.SaveRequest_Credentials{Credentials: &pb.Credentials{Login: login, Password: password}},
		}
	case "4":
		say("Enter binary")
		s.Scan()
		b := s.Bytes()
		in = &pb.SaveRequest{
			Kind: pb.Kind_BINARY,
			Data: &pb.SaveRequest_Binary{Binary: b},
		}
	case "0":
		say("Returning to previous menu..")
		time.Sleep(500 * time.Millisecond)
		return nil
	default:
		say("Wrong value, returning to previous menu..")
		time.Sleep(500 * time.Millisecond)
		return nil
	}
	say("Metadata (not necessary)")
	s.Scan()
	in.Meta = s.Text()

	return in
}

// SyncData method gets all user's data from server and saves it into local cache.
func (c *Client) SyncData(ctx context.Context) error {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("token", c.token))

	data, err := c.grpcClient.GetMany(ctx, &pb.GetManyRequest{})
	if err != nil {
		return err
	}

	cnt := 0
	for _, d := range data.UserData {
		c.localCache.Store(d.ID, d)
		cnt += 1
	}

	sayf("\nCurrent user: %d record(s) found\n", cnt)

	return nil
}

func (c *Client) ListenAction(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		say(InputAction)
		scanner.Scan()
		input := scanner.Text()
		switch strings.TrimSpace(input) {
		case "1":
			if _, err := c.Register(ctx, scanner); err != nil {
				sayf("failed to register user: %v\n", err)
			}

		case "2":
			if _, err := c.LogIn(ctx, scanner); err != nil {
				if IsServerError(err) {
					say(ErrServerUnavailable.Error())
					continue
				}
				sayf("failed to log in user: %v\n", err)
			}

			if err := c.SyncData(ctx); err != nil {
				if IsServerError(err) {
					say(ErrServerUnavailable.Error())
					continue
				}
				sayf("Failed to sync data with server, err=%+v", err.Error())
			} else {
				say("Remote data successfully synchronized..")
			}

		case "3":
			got, err := c.GetData(ctx, scanner)
			if err != nil {
				if IsServerError(err) {
					say(ErrServerUnavailable.Error())
					// если серверная ошибка, идем рыться в локальную мапу
					//c.localCache.Range(func(key, value interface{}) bool {
					//	fmt.Printf("Key: %s, Value: %s\n", key.(string), value)
					//	return true
					//})
					continue
				}
				sayf("\nfailed to get data: %v\n", err)
				continue
			}
			say(got)
		case "4":
			if err := c.SaveData(ctx, scanner); err != nil {
				if IsServerError(err) {
					say(ErrServerUnavailable.Error())
					continue
				}
				sayf("failed to save data: %v\n", err)
			}
		case "0":
			say("Exiting the application..")
			os.Exit(1)
		default:
			say("Unknown value")
		}
	}
}

// IsServerError checks for server unavailability and returns bool.
func IsServerError(err error) bool {
	serverErrors := []codes.Code{codes.Unavailable}
	if e, ok := status.FromError(err); ok {
		if !slices.Contains(serverErrors, e.Code()) {
			return false
		}
	}
	return true
}

var (
	ErrUserUnauthenticated = errors.New("unauthenticated user")
	ErrServerUnavailable   = errors.New("server unavailable")
)
