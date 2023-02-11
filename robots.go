package robots

import (
	"net/url"
	"regexp"
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
			} else {
				disableUA = true
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
	allow := true

	for _, v := range robots.Disallowlist {
		v = strings.ReplaceAll(v, "*", ".*?")
		/* check prefix */
		if v[len(v)-1] != '$' || v[len(v)-1] == '/' {
			v += ".*"
		}
		r, err := regexp.Compile(v)
		/* if error, skip it. */
		if err == nil {
			if r.Match([]byte(target.String())) {
				allow = false
			}
		}
	}

	for _, v := range robots.Allowlist {
		v = strings.ReplaceAll(v, "*", ".*?")
		/* check prefix */
		if v[len(v)-1] != '$' || v[len(v)-1] == '/' {
			v += ".*"
		}
		r, err := regexp.Compile(v)
		/* if error, skip it. */
		if err == nil {
			if r.Match([]byte(target.String())) {
				allow = true
			}
		}
	}

	return allow
}
