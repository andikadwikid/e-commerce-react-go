package middlewares

import (
	"fmt"
	"strings"

	"backend-commerce/config"
)

func GetCORSOrigins() []string {
	originsEnv := config.GetEnv("FRONTEND_URL", "")
	rawOrigins := strings.Split(originsEnv, ",")
	uniqueMap := make(map[string]bool)
	fmt.Println("FRONTEND_URL:", originsEnv)
	fmt.Println("rawOrigins:", rawOrigins)
	var cleanOrigins []string
	for _, o := range rawOrigins {
		trimmed := strings.TrimSpace(o)
		trimmed = strings.Trim(trimmed, "\"'")
		trimmed = strings.TrimRight(trimmed, "/")

		if trimmed == "" {
			continue
		}

		if !uniqueMap[trimmed] {
			uniqueMap[trimmed] = true
			cleanOrigins = append(cleanOrigins, trimmed)
		}
	}

	return cleanOrigins
}
