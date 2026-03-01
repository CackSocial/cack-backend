package mentions

import "regexp"

var mentionRegex = regexp.MustCompile(`(?:^|[^a-zA-Z0-9_])@([a-zA-Z0-9_]+)`)

// ExtractMentions returns a deduplicated list of usernames mentioned in the text.
// Mentions are in the format @username.
func ExtractMentions(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]struct{})
	var mentions []string
	for _, match := range matches {
		username := match[1]
		if _, ok := seen[username]; !ok {
			seen[username] = struct{}{}
			mentions = append(mentions, username)
		}
	}
	return mentions
}
