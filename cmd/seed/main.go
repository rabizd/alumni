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
	jobPostings := flag.Int("jobs", 10, "random job postings to add, each by a random existing user")
	fresh := flag.Bool("fresh", false, "empty every table first")
	randSeed := flag.Int64("seed", time.Now().UnixNano(), "random seed; the same number gives the same records")
	admin := flag.String("admin", "", "make the user with this email an admin, then stop (adds nothing)")
	flag.Parse()

	ctx := context.Background()
	if err := model.Connect(ctx); err != nil {
		log.Fatalf("database: %v", err)
	}

	// Only an admin can change roles on the site, so the first admin is made here.
	if *admin != "" {
		u, _, found, err := model.Users.FindForLogin(*admin)
		if err != nil {
			log.Fatalf("admin: %v", err)
		}
		if !found {
			log.Fatalf("admin: no user with the email %q; register on /register first", *admin)
		}
		if _, _, err := model.Users.SetRole(u.ID, model.RoleAdmin); err != nil {
			log.Fatalf("admin: %v", err)
		}
		log.Printf("%s (%s, id %d) is now an admin", u.Name, u.Email, u.ID)
		return
	}

	if *fresh {
		if err := model.Truncate(ctx); err != nil {
			log.Fatalf("fresh: %v", err)
		}
		log.Print("emptied every table")
	}

	counts := seed.Counts{Users: *users, Announcements: *announcements, Alumni: *alumni, JobPostings: *jobPostings}
	if err := seed.Run(seed.NewFactory(*randSeed), counts); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Printf("added %d random users, %d announcements, %d alumni, %d job postings, plus the samples in any empty table (seed %d)",
		*users, *announcements, *alumni, *jobPostings, *randSeed)
}
