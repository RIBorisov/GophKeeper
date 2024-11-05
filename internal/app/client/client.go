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

func NewClient() (*Client, error) {
	creds, err := applyTLS()
	if err != nil {
		return nil, fmt.Errorf("failed to apply tls: %w", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	return &Client{
		grpcClient: pb.NewGophKeeperServiceClient(conn),
		localCache: &sync.Map{},
		token:      "",
	}, nil
}

func scanCredentials(s *bufio.Scanner) (string, string) {
	say("Enter login")
	s.Scan()
	login := s.Text()
	say("Enter password")
	s.Scan()
	password := s.Text()
	return login, password
}

func (c *Client) Register(ctx context.Context, s *bufio.Scanner) (*pb.RegisterResponse, error) {
	login, password := scanCredentials(s)
	in := &pb.RegisterRequest{Login: login, Password: password}
	resp, err := c.grpcClient.Register(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	c.token = resp.GetToken()

	return resp, nil
}

func (c *Client) LogIn(ctx context.Context, s *bufio.Scanner) (*pb.AuthResponse, error) {
	login, password := scanCredentials(s)
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
	input := strings.TrimSpace(s.Text())

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("token", c.token))
	resp, err := c.grpcClient.Get(ctx, &pb.GetRequest{ID: input})
	if err != nil {
		if IsServerError(err) {
			say(ErrServerUnavailable.Error())
			say("going to get data from local storage..")
			v, ok := c.localCache.Load(input)
			if !ok {
				return nil, err
			}
			local, ok := v.(*pb.GetResponse)
			if !ok {
				say("found but can't recognize value, possible corrupted data")
				return nil, errors.New("can't recognize data")
			}

			return local, nil
		}

		return nil, err
	}

	if resp == nil {
		return nil, ErrDataNotFound
	}

	return resp, nil
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

		in.Kind = pb.Kind_CARD
		in.Data = &pb.SaveRequest_Card{Card: &pb.Card{Number: number, MmYy: mmYY, Cvc: cvc}}
	case "2":
		say("Enter text")
		s.Scan()
		text := s.Text()

		in.Kind = pb.Kind_TEXT
		in.Data = &pb.SaveRequest_Text{Text: text}
	case "3":
		say("Enter login")
		s.Scan()
		login := s.Text()
		say("Enter password")
		s.Scan()
		password := s.Text()

		in.Kind = pb.Kind_CREDENTIALS
		in.Data = &pb.SaveRequest_Credentials{Credentials: &pb.Credentials{Login: login, Password: password}}
	case "4":
		say("Enter binary")
		s.Scan()
		b := s.Bytes()

		in.Kind = pb.Kind_BINARY
		in.Data = &pb.SaveRequest_Binary{Binary: b}
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
	for _, d := range data.GetUserData() {
		c.localCache.Store(d.GetID(), d)
		cnt++
	}

	sayf("\nCurrent user: %d record(s) found\n", cnt)

	return nil
}

func (c *Client) ListenAction(ctx context.Context, bDate, bVersion string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		say(InputAction)
		scanner.Scan()
		input := scanner.Text()
		switch strings.TrimSpace(input) {
		case Register.String():
			if _, err := c.Register(ctx, scanner); err != nil {
				sayf("failed to register user: %v\n", err)
			}

		case LogIn.String():
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

		case GetData.String():
			got, err := c.GetData(ctx, scanner)
			if err != nil {
				sayf("failed to get data: %v\n", err)
				continue
			}
			say("Result: " + got.String())
		case SaveData.String():
			if err := c.SaveData(ctx, scanner); err != nil {
				if IsServerError(err) {
					say(ErrServerUnavailable.Error())
					continue
				}
				sayf("failed to save data: %v\n", err)
			}
		case Return.String():
			say("Exiting the application..")
			os.Exit(1)
		case BuildInfo.String():
			sayf("Build date: %s\n", bDate)
			sayf("Build version: %s\n", bVersion)

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
	ErrDataNotFound        = errors.New("not found data")
)
