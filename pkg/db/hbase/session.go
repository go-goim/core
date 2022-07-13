package hbase

import (
	"context"
	"fmt"

	"github.com/tsuna/gohbase/hrpc"
)

type session struct {
	ctx               context.Context
	table             string
	key               string
	family            *string
	qualifier         *string
	amount            *int64
	startRow, stopRow *string
	values            map[string]map[string][]byte
	expectedValue     []byte
	// options
	opts []func(hrpc.Call) error
}

func newSession(ctx context.Context) *session {
	return &session{
		ctx: ctx,
		// options
		opts: make([]func(hrpc.Call) error, 0),
	}
}

func (s *session) WithContext(ctx context.Context) *session {
	s.ctx = ctx
	return s
}

func (s *session) WithTable(table string) *session {
	s.table = table
	return s
}

func (s *session) WithKey(key string) *session {
	s.key = key
	return s
}

func (s *session) WithFamily(family string) *session {
	s.family = &family
	return s
}

func (s *session) WithQualifier(qualifier string) *session {
	s.qualifier = &qualifier
	return s
}

func (s *session) WithAmount(amount int64) *session {
	s.amount = &amount
	return s
}

func (s *session) WithRange(startRow, stopRow string) *session {
	s.startRow = &startRow
	s.stopRow = &stopRow
	return s
}

func (s *session) WithValues(values map[string]map[string][]byte) *session {
	s.values = values
	return s
}

func (s *session) WithExpectedValue(expectedValue []byte) *session {
	s.expectedValue = expectedValue
	return s
}

func (s *session) WithOptions(opts ...func(hrpc.Call) error) *session {
	s.opts = opts
	return s
}

func (s *session) isSetRange() bool {
	return s.startRow != nil && s.stopRow != nil
}

func (s *session) isSetFamily() bool { // nolint: unused
	return s.family != nil
}

func (s *session) isSetQualifier() bool { // nolint: unused
	return s.qualifier != nil
}

func (s *session) isSetAmount() bool {
	return s.amount != nil
}

var (
	ErrNilContext = fmt.Errorf("nil context")
	ErrNilTable   = fmt.Errorf("nil table")
)

// TODO: more validation for different operations.
func (s *session) validate() error {
	if s.ctx == nil {
		return ErrNilContext
	}
	if s.table == "" {
		return ErrNilTable
	}

	return nil
}

type Result interface {
	Scanner() hrpc.Scanner
	Result() *hrpc.Result
	Int64() int64
	Bool() bool
	Err() error
}

var _ Result = &result{}

type result struct {
	result  *hrpc.Result
	scanner hrpc.Scanner
	i64     int64 // for Increment
	b       bool  // for CheckAndPut
	err     error
}

func (r *result) Scanner() hrpc.Scanner {
	return r.scanner
}

func (r *result) Result() *hrpc.Result {
	return r.result
}

func (r *result) Int64() int64 {
	return r.i64
}

func (r *result) Bool() bool {
	return r.b
}

func (r *result) Err() error {
	return r.err
}
