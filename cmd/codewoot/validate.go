package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/sdrvirtual/codewoot/internal/codechat"
// )

// func ValidateCodechatSendText(ctx context.Context, baseURL, globalToken, instanceToken, instanceName, toNumber, text string) error {
// 	c, err := codechat.New(baseURL, globalToken, codechat.WithInstanceToken(instanceToken, instanceName))
// 	if err != nil {
// 		return err
// 	}

// 	cu := map[string]any{
// 		"text": text,
// 	}
// 	payload := map[string]any{
// 		"number": toNumber,
// 		"textMessage": cu,
// 	}

// 	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
// 		var cancel context.CancelFunc
// 		ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
// 		defer cancel()
// 	}

// 	res, err := c.SendText(ctx, instanceName, payload)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("Codechat SendText response: %s", string(res))
// 	return nil
// }
