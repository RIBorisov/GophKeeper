package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

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
	token      string
}

func NewClient(ctx context.Context) (*Client, error) {
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	return &Client{grpcClient: pb.NewGophKeeperServiceClient(conn)}, nil
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
	}
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
				sayf("failed to log in user: %v\n", err)
			}
		case "3":
			got, err := c.GetData(ctx, scanner)
			if err != nil {
				sayf("failed to get data: %v\n", err)
			}
			say(got)
		case "4":
			if err := c.SaveData(ctx, scanner); err != nil {
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
