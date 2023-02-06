package robots

import (
	"net/url"
	"strings"

	"github.com/HonokaNo/cacheget"
)

type Robots struct {
	Allowlist    []string
	Disallowlist []string
	Sitemaps     []string
}

func Parse(url, UA string) (Robots, error) {
	var robots Robots

	bbody, _, err := cacheget.CacheGet(url)
	if err != nil {
		return Robots{[]string{}, []string{}, []string{}}, err
	}
	body := string(bbody)

	body = strings.ReplaceAll(body, "\r\n", "\n")

	foundUA := false
	disableUA := true

	for _, v := range strings.Split(body, "\n") {
		if strings.HasPrefix(v, "User-agent: ") {
			arg := v[len("User-agent: "):]

			if arg == UA {
				foundUA = true
				disableUA = false
			} else if !foundUA && arg == "*" {
				disableUA = false
			}
		} else if !disableUA {
			if strings.HasPrefix(v, "Allow: ") {
				robots.Allowlist = append(robots.Allowlist, v[len("Allow: "):])
			} else if strings.HasPrefix(v, "Disallow: ") {
				robots.Disallowlist = append(robots.Disallowlist, v[len("Disallow: "):])
			}
		}

		if strings.HasPrefix(v, "Sitemap: ") {
			robots.Sitemaps = append(robots.Sitemaps, v[len("Sitemap: "):])
		}
	}

	return robots, nil
}

func IsAllowURL(target url.URL, robots Robots) bool {
	search := target
	allow := true

	for _, v := range robots.Disallowlist {
		search.Path = v[1:]

		if strings.HasPrefix(target.String(), search.String()) {
			allow = false
			break
		}
	}

	for _, v := range robots.Allowlist {
		search.Path = v[1:]

		if !allow && strings.HasPrefix(target.String(), search.String()) {
			allow = true
			break
		}
	}

	return allow
}
