package mentions

import (
	"reflect"
	"testing"
)

func TestExtractMentions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "empty string",
			content: "",
			want:    nil,
		},
		{
			name:    "single mention",
			content: "@alice hello",
			want:    []string{"alice"},
		},
		{
			name:    "multiple mentions",
			content: "@alice hey @bob",
			want:    []string{"alice", "bob"},
		},
		{
			name:    "duplicate mentions",
			content: "@alice @alice",
			want:    []string{"alice"},
		},
		{
			name:    "mixed with hashtags",
			content: "#tag @user",
			want:    []string{"user"},
		},
		{
			name:    "no mentions",
			content: "hello world",
			want:    nil,
		},
		{
			name:    "mentions with underscores",
			content: "@user_name",
			want:    []string{"user_name"},
		},
		{
			name:    "adjacent mentions",
			content: "@alice@bob",
			want:    []string{"alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractMentions(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractMentions(%q) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}
