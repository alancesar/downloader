package param

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	Recent    Order = "recent"
	Followers Order = "followers"
	Best      Order = "best"
	Trending  Order = "trending"
	Top28     Order = "top28"
)

type (
	Order string

	Search interface {
		ToURL() string
	}

	AuthorSearch struct {
		Username string
		Page     int
		Order    Order
		Verified bool
		Tags     []string
	}
)

func (p AuthorSearch) ToURL() string {
	rawURL := fmt.Sprintf("/users/%s/search", p.Username)
	parsedURL, _ := url.Parse(rawURL)
	query := parsedURL.Query()

	if p.Page != 0 {
		query.Add("page", strconv.Itoa(p.Page))
	}

	if p.Order != "" {
		query.Add("order", string(p.Order))
	}

	if p.Verified {
		query.Add("verified", "y")
	}

	if len(p.Tags) > 0 {
		query.Add("tags", strings.Join(p.Tags, ","))
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}
