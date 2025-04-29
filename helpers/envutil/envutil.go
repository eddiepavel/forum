package envutil

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadEnv(filename string) error {

	file, error := os.Open(filename)

	if error != nil {
		return error
	}

	defer file.Close()

	error = scanEnvFile(file)

	return error
}

func scanEnvFile(file *os.File) error {

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split at the first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Skipping invalid line in .env: %s", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if any
		val = strings.Trim(val, `"'`)

		os.Setenv(key, val)
	}

	return scanner.Err()
}

func GetEnvString(key string) string {
	val, exists := os.LookupEnv(key)

	if !exists {
		return ""
	}

	return val
}
