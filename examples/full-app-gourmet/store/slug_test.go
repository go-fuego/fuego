package store

import "testing"

func TestSlug(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "simple",
			want: "simple",
		},
		{
			name: "simple with space",
			want: "simple-with-space",
		},
		{
			name: "simple with space and accént",
			want: "simple-with-space-and-accent",
		},
		{
			name: "simple with space and accent and ✅",
			want: "simple-with-space-and-accent-and-%E2%9C%85",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slug(tt.name); got != tt.want {
				t.Errorf("slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzSlug(f *testing.F) {
	f.Add("simple")
	f.Add("simple with space")
	f.Add("simple with space and accént")
	f.Add("simple with space and accent and ✅")
	f.Fuzz(func(t *testing.T, name string) {
		slug(name)
	})
}
