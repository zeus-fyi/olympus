package hermes_email_notifications

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

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

type EmailContents struct {
	From    string
	Subject string
	Body    string
}

func (g *GmailServiceClient) GetReadEmails(email string) ([]EmailContents, error) {
	r, err := g.Users.Messages.List(email).MaxResults(5).Do()
	if err != nil {
		return nil, err
	}
	var emails []EmailContents // Slice of EmailContents instead of gmail.Message
	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
		return emails, nil
	}
	fmt.Println("Messages:")
	for _, m := range r.Messages {
		msg, err := g.Users.Messages.Get(email, m.Id).Format("full").Do()
		if err != nil {
			return nil, err
		}

		var emailContents EmailContents
		if msg != nil {
			// Extracting the headers for sender and subject
			for _, header := range msg.Payload.Headers {
				if header.Name == "From" {
					// Use regular expression to extract just the email address
					re := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
					emailMatches := re.FindStringSubmatch(header.Value)
					if len(emailMatches) > 0 {
						emailContents.From = emailMatches[0] // Assign the extracted email address
					}
				}
				if header.Name == "Subject" {
					emailContents.Subject = header.Value
				}
			}

			// Extracting the body of the email
			body := ""
			if msg.Payload.Parts == nil {
				body = decodeBase64URL(msg.Payload.Body.Data)
			} else {
				for _, part := range msg.Payload.Parts {
					if part.MimeType == "text/plain" {
						body += decodeBase64URL(part.Body.Data)
					} else if part.MimeType == "multipart/alternative" {
						for _, subPart := range part.Parts {
							if subPart.MimeType == "text/plain" {
								body += decodeBase64URL(subPart.Body.Data)
							}
						}
					}
				}
			}
			emailContents.Body = body // Assign the decoded body to the struct

			emails = append(emails, emailContents) // Append the constructed EmailContents to the slice
		}
	}
	return emails, nil
}

// Helper function to decode base64 URL encoded strings
func decodeBase64URL(base64Message string) string {
	data, err := base64.URLEncoding.DecodeString(base64Message)
	if err != nil {
		fmt.Println("Error decoding base64 message: ", err)
		return ""
	}
	return string(data)
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
