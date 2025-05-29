package flash

import (
	"net/http"
	"net/url"
	"strings"
)

// HandleMessages handles flash messages by appending them to the query parameters of a redirect URL.
func HandleMessages(w http.ResponseWriter, r *http.Request, messages map[string]string, redirectURL string, messageType string) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		http.Error(w, "Invalid redirect URL", http.StatusInternalServerError)
		return
	}

	var info []string
	q := u.Query()

	for _, values := range messages {
		if len(values) > 0 {
			info = append(info, values)
		}
	}

	q.Set(messageType, strings.Join(info, ","))

	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
