package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	messagev1 "github.com/go-goim/api/message/v1"
	"github.com/go-goim/core/pkg/types"
)

//nolint:scopelint
func TestSession(t *testing.T) {
	type args struct {
		tp   messagev1.SessionType
		from types.ID
		to   types.ID
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
				from: 71990933687635970,
				to:   71990933687635971,
			},
			want: "000aG9PKEB8ch0aG9PKEB8ci",
		},
		{
			name: "group chat",
			args: args{
				tp:   1,
				from: 71990933687635972,
				to:   71990933687635973,
			},
			want: "010aG9PKEB8cj00000000000",
		},
		{
			name: "broadcast",
			args: args{
				tp:   2,
				from: 71990933687635970,
				to:   71990933687635970,
			},
			want: "020000000000000000000000",
		},
		{
			name: "channel",
			args: args{
				tp:   3,
				from: 71990933687635974,
				to:   71990933687635975,
			},
			want: "030aG9PKEB8cm0aG9PKEB8cn",
		},
		{
			name: "feature types",
			args: args{
				tp:   127,
				from: 71990933687635976,
				to:   71990933687635977,
			},
			want: "7f0aG9PKEB8co0aG9PKEB8cp",
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

//nolint:scopelint
func TestParseSession(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantTp   int32
		wantFrom types.ID
		wantTo   types.ID
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name: "single chat",
			args: args{
				s: "000aG9PKEB8ch0aG9PKEB8ci",
			},
			wantTp:   0,
			wantFrom: types.ID(71990933687635970),
			wantTo:   types.ID(71990933687635971),
			wantErr:  assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTp, gotFrom, gotTo, err := ParseSession(tt.args.s) //nolint:scopelint
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSession(%v)", tt.args.s)) {
				return
			}
			assert.Equalf(t, tt.wantTp, gotTp, "ParseSession(%v)", tt.args.s)
			assert.Equalf(t, tt.wantFrom, gotFrom, "ParseSession(%v)", tt.args.s)
			assert.Equalf(t, tt.wantTo, gotTo, "ParseSession(%v)", tt.args.s)
		})
	}
}
