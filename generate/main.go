package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

const key = "key"

func main() {
	const times = 10_000

	// Setup the value we want to retrieve in each iteration
	ctx := context.WithValue(context.Background(), key, "some-val")
	fat(ctx, times)

	ctx = context.WithValue(context.Background(), key, "some-val")
	thin(ctx, times)

	ctx = context.WithValue(context.Background(), key, "some-val")
	experiment(ctx, times)
}

func fat(ctx context.Context, times uint64) {
	for range times {
		// wrap the context, each iteration makes it bigger
		ctx = contextWithRandom(ctx)

		start := time.Now()
		// simulate the logging lib which retrieves context values
		_ = ctx.Value(key)

		fmt.Printf("fat,%v\n", time.Since(start).Nanoseconds())
	}
}

func thin(ctx context.Context, times uint64) {
	for range times {
		// shadow the context, each iteration creates a new one and it doesn't grow
		ctx := contextWithRandom(ctx)

		start := time.Now()
		_ = ctx.Value(key)

		fmt.Printf("thin,%v\n", time.Since(start).Nanoseconds())
	}
}

func experiment(ctx context.Context, times uint64) {
	wrapper := something()
	r := &R{ctx}
	for range times {
		wrapper(r)
		start := time.Now()
		_ = r.Ctx.Value(key)
		fmt.Printf("experiment,%v\n", time.Since(start).Nanoseconds())
	}
}

func contextWithRandom(ctx context.Context) context.Context {
	return context.WithValue(ctx, "other_key", uuid.Must(uuid.NewV4()))
}

type R struct {
	Ctx context.Context
}

func something() func(*R) {
	return func(r *R) {
		ctx := r.Ctx
		ctx = contextWithRandom(ctx)
		r.Ctx = ctx // triggered on this line
	}
}
