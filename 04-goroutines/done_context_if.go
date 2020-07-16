package main

import (
	"context"
	"fmt"
)

type something struct{}

// start_main OMIT
func main() {
	ctx, cancel := context.WithCancel(context.Background()) // HL
	fmt.Println("ctx.Err:", ctx.Err())
	cancel()

	err := slowComputation(ctx)
	fmt.Printf("slow computation fished:\n\t%v\n", err)
}

// end_main OMIT

// start_slowComputation OMIT
func slowComputation(ctx context.Context) error {
	if err := ctx.Err(); err != nil { // Check if we should proceed // HL
		return fmt.Errorf("context already done: %w", err)
	}

	// Something slow...

	fmt.Printf("successfuly fineshed slowComputation")
	return nil
}

// end_slowComputation OMIT
