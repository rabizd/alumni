// Command seed fills the database: the fixed samples plus random records from
// the factories. Run it inside Docker with `docker compose run --rm seed`.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/seed"
)

func main() {
	users := flag.Int("users", 10, "random users to add")
	announcements := flag.Int("announcements", 5, "random announcements to add")
	alumni := flag.Int("alumni", 10, "random graduates to add")
	fresh := flag.Bool("fresh", false, "empty every table first")
	randSeed := flag.Int64("seed", time.Now().UnixNano(), "random seed; the same number gives the same records")
	flag.Parse()

	ctx := context.Background()
	if err := model.Connect(ctx); err != nil {
		log.Fatalf("database: %v", err)
	}

	if *fresh {
		if err := model.Truncate(ctx); err != nil {
			log.Fatalf("fresh: %v", err)
		}
		log.Print("emptied every table")
	}

	counts := seed.Counts{Users: *users, Announcements: *announcements, Alumni: *alumni}
	if err := seed.Run(seed.NewFactory(*randSeed), counts); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Printf("added %d random users, %d announcements, %d alumni, plus the samples in any empty table (seed %d)",
		*users, *announcements, *alumni, *randSeed)
}
