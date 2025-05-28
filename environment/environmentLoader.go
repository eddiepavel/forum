package environment

import (
	"bufio"
	"os"
	"strings"
)

// loadDotEnv reads key=value lines from the given file path
// and sets them into the process’s environment via os.Setenv.
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		// If the file doesn’t exist, just skip without error
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		// skip comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// split key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		os.Setenv(key, val)
	}
	return scanner.Err()
}
