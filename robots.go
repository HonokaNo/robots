package robots

import (
	"io"
	"net/http"
	"strings"
)

type Robots struct {
	Allowlist    []string
	Disallowlist []string
	Sitemaps     []string
}

func Parse(url, UA string) (Robots, error) {
	var robots Robots

	resp, err := http.Get(url)
	if err != nil {
		return robots, err
	}

	bbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return robots, err
	}
	body := string(bbody)

	body = strings.ReplaceAll(body, "\r\n", "\n")

	disableUA := false

	for _, v := range strings.Split(body, "\n") {
		if strings.HasPrefix(v, "User-agent: ") {
			arg := v[len("User-agent: "):]

			disableUA = !(arg == "*" || arg == UA)
		} else if !disableUA {
			if strings.HasPrefix(v, "Allow: ") {
				robots.Allowlist = append(robots.Allowlist, v[len("Allow: "):])
			} else if strings.HasPrefix(v, "Disallow: ") {
				robots.Disallowlist = append(robots.Disallowlist, v[len("Disallow: "):])
			} else if strings.HasPrefix(v, "Sitemap: ") {
				robots.Sitemaps = append(robots.Sitemaps, v[len("Sitemap: "):])
			}
		}
	}

	return robots, nil
}
