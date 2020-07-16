package main

import (
	"context"
	"fmt"
	"time"
)

type something struct{}

// start_main OMIT
func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second) // HL
	err := slowComputation(ctx)
	fmt.Printf("slow computation fished:\n\t%v\n", err)
}

// end_main OMIT

// start_slowComputation OMIT
func slowComputation(ctx context.Context) error {
	ready := make(chan something)
	go func() { time.Sleep(time.Minute); ready <- something{} }()

	select {
	case <-ctx.Done(): // Will receive a message when the context times out // HL
		return fmt.Errorf("slow computation was too slow: %w", ctx.Err())
	case value := <-ready:
		fmt.Printf("successfuly fineshed: %#v\n", value)
		return nil
	}
}

// end_slowComputation OMIT
