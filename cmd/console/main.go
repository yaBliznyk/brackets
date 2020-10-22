package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/yaBliznyk/brackets/internal/brackets"
	"os"
)

// Main function
// Run service and listen for errors
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	args := os.Args
	var str string
	if len(args) == 1 {
		return errors.New("необходимо передать строку для проверки и корректировки")
	} else {
		str = args[1]
	}

	err, bracketSvc := brackets.NewBaseService(&brackets.Config{
		Logger: nil,
		Bkts:   "{}[]()",
	})
	if err != nil {
		return err
	}

	ctx := context.Background()

	v := bracketSvc.Validate(ctx, str)
	err, fix := bracketSvc.Fix(ctx, str)
	if err != nil {
		return err
	}

	fmt.Println("Is valid:", v, "\nFixed as:", fix)

	return nil
}
