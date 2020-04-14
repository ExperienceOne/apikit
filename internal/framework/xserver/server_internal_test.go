package xserver

import (
	"testing"

	"github.com/go-ozzo/ozzo-routing"
)

func TestNewServer(t *testing.T) {

	tests := []struct {
		name       string
		args       []Middleware
		wantBefore bool
		wantAfter  bool
	}{
		{
			name: "NoMiddleware",
			args: nil,
		},
		{
			name: "OnlyBeforeMiddleware",
			args: []Middleware{
				{
					Handler: func(c *routing.Context) error { return nil },
					After:   false,
				},
			},
			wantBefore: true,
		},
		{
			name: "OnlyAfterMiddleware",
			args: []Middleware{
				{
					Handler: func(c *routing.Context) error { return nil },
					After:   true,
				},
			},
			wantAfter: true,
		},
		{
			name: "BeforeAndAfterMiddleware",
			args: []Middleware{
				{
					Handler: func(c *routing.Context) error { return nil },
					After:   false,
				},
				{
					Handler: func(c *routing.Context) error { return nil },
					After:   true,
				},
			},
			wantAfter:  true,
			wantBefore: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got := NewServer(&ServerOpts{Middleware: tc.args})
			if (len(got.after) != 0) != tc.wantAfter {
				t.Fatalf("xserver.NewServer().after = %v, want %v", got.after, tc.wantAfter)
			}
			if (len(got.before) != 0) != tc.wantBefore {
				t.Fatalf("xserver.NewServer().before = %v, want %v", got.before, tc.wantAfter)
			}
		})
	}
}
