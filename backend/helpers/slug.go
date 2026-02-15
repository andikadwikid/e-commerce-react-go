package helpers

import (
	"regexp"
	"strings"

)

func Slugify(text string) string {
	slug := strings.TrimSpace(text)

	slug = strings.ToLower(slug)

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	slug = strings.ReplaceAll(slug, " ", "-")

	reDoubleDash := regexp.MustCompile(`-+`)
	slug = reDoubleDash.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
