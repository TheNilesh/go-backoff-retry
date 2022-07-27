package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	backoff "github.com/cenkalti/backoff/v4"
)

var Err4xx = errors.New("404")
var Err5xx = errors.New("500")

func main() {

	/*
		Use case where we should not retry if error is client side error.
		This can be achieved using PermanentErr
	*/
	doSomething := func(i int) error {
		switch i {
		case 0:
			fmt.Println("4xx")
			return Err4xx
		case 1:
			fmt.Println("No error")
			return nil
		case 2:
			fmt.Println("5xx")
			return Err5xx
		default:
			fmt.Println("connection error")
			return errors.New("connection refused")
		}
	}

	ctx := context.Background()
	err := backoff.Retry(func() error {
		err := doSomething(rand.Intn(3))
		if errors.Is(err, Err4xx) {
			return backoff.Permanent(err) // This wont be retried
		}
		return err
	}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		fmt.Println("err", err)
		return
	}
}
