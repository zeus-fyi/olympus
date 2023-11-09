package hermes_email_notifications

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	MainEmailUser    GmailServiceClient
	SupportEmailUser GmailServiceClient
)

type GmailServiceClient struct {
	*gmail.Service
}

func (g *GmailServiceClient) ReadEmails(email string) {
	r, err := g.Users.Messages.List(email).MaxResults(5).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
	} else {
		fmt.Println("Messages:")
		for _, m := range r.Messages {
			msg, err := g.Users.Messages.Get(email, m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}
			fmt.Printf("- %v (snippet: '%v')\n", m.Id, msg.Snippet)
		}
	}
}

func InitNewGmailServiceClients(ctx context.Context, authJsonBytes []byte) {
	MainEmailUser = NewGmailServiceClient(ctx, authJsonBytes, "alex@zeus.fyi")
	SupportEmailUser = NewGmailServiceClient(ctx, authJsonBytes, "support@zeus.fyi")
	return
}

func NewGmailServiceClient(ctx context.Context, authJsonBytes []byte, email string) GmailServiceClient {
	// Read the service account key file
	// Authenticate and create the service

	conf, err := google.JWTConfigFromJSON(authJsonBytes, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("JWTConfigFromJSON: %v", err)
	}
	conf.Subject = email
	client := conf.Client(ctx)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return GmailServiceClient{srv}
}

func NewGmail(ctx context.Context, authJsonBytes []byte, email string) {
	// Read the service account key file
	// Authenticate and create the service
	conf, err := google.JWTConfigFromJSON(authJsonBytes, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("JWTConfigFromJSON: %v", err)
	}
	conf.Subject = email
	client := conf.Client(ctx)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Get the last 5 messages
	user := email
	r, err := srv.Users.Messages.List(user).MaxResults(5).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
	} else {
		fmt.Println("Messages:")
		for _, m := range r.Messages {
			msg, err := srv.Users.Messages.Get(user, m.Id).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
			}
			fmt.Printf("- %v (snippet: '%v')\n", m.Id, msg.Snippet)
		}
	}
}
