package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	type args struct {
		tp   int32
		from int64
		to   int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single chat",
			args: args{
				tp:   0,
				from: 1,
				to:   2,
			},
			want: "0000000000000000000000100000000000000000002",
		},
		{
			name: "group chat",
			args: args{
				tp:   1,
				from: 1,
				to:   2,
			},
			want: "0010000000000000000000200000000000000000000",
		},
		{
			name: "broadcast",
			args: args{
				tp:   2,
				from: 1,
				to:   2,
			},
			want: "0020000000000000000000000000000000000000000",
		},
		{
			name: "channel",
			args: args{
				tp:   3,
				from: 1,
				to:   2,
			},
			want: "0030000000000000000000100000000000000000002",
		},
		{
			name: "invalid type -1",
			args: args{
				tp:   -1,
				from: 1,
				to:   2,
			},
			want: "",
		},
		{
			name: "invalid type 256",
			args: args{
				tp:   256,
				from: 1,
				to:   2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Session(tt.args.tp, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("Session() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSession(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantTp   int32
		wantFrom int64
		wantTo   int64
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "single chat",
			args: args{
				s: "0000000000000000000000100000000000000000002",
			},
			wantTp:   0,
			wantFrom: 1,
			wantTo:   2,
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTp, gotFrom, gotTo, err := ParseSession(tt.args.s)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSession(%v)", tt.args.s)) {
				return
			}
			assert.Equalf(t, tt.wantTp, gotTp, "ParseSession(%v)", tt.args.s)
			assert.Equalf(t, tt.wantFrom, gotFrom, "ParseSession(%v)", tt.args.s)
			assert.Equalf(t, tt.wantTo, gotTo, "ParseSession(%v)", tt.args.s)
		})
	}
}
