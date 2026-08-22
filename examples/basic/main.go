// Command basic demonstrates listing and fetching matches with the WFA API client.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/client"
	"github.com/TheWFA/api-client-go/matches"
)

func main() {
	c, err := client.New(os.Getenv("WFA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	page, err := c.Matches.List(ctx, matches.ListQuery{
		OrderByDateDesc: wfa.Bool(true),
		ListParams:      wfa.ListParams{ItemsPerPage: wfa.Int(5)},
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(page.Items) == 0 {
		fmt.Println("no matches found")
		return
	}

	for _, m := range page.Items {
		fmt.Printf("#%d: %s %d-%d %s\n", m.ID, m.HomeTeam.Name, m.HomeScore, m.AwayScore, m.AwayTeam.Name)
	}

	full, err := c.Matches.Get(ctx, page.Items[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("first match had %d events\n", len(full.Events))
}
