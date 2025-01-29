//package main

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

func main() {
	apiKey := "re_QUQsC8qH_2rznRRpCysRkGNT4npLurqge"

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{"vanspaul09@gmail.com"},
		Subject: "Hello World",
		Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(sent)
	}
}
