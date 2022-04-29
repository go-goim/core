// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/config/v1/config.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Server with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Server) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Server with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ServerMultiError, or nil if none found.
func (m *Server) ValidateAll() error {
	return m.validate(true)
}

func (m *Server) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if _, ok := _Server_Scheme_InLookup[m.GetScheme()]; !ok {
		err := ServerValidationError{
			field:  "Scheme",
			reason: "value must be in list [http grpc tcp]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetAddr() != "" {

		if err := m._validateHostname(m.GetAddr()); err != nil {
			if ip := net.ParseIP(m.GetAddr()); ip == nil {
				err := ServerValidationError{
					field:  "Addr",
					reason: "value must be a valid hostname, or ip address",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}
		}

	}

	if val := m.GetPort(); val <= 10000 || val >= 60535 {
		err := ServerValidationError{
			field:  "Port",
			reason: "value must be inside range (10000, 60535)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ServerMultiError(errors)
	}

	return nil
}

func (m *Server) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

// ServerMultiError is an error wrapping multiple validation errors returned by
// Server.ValidateAll() if the designated constraints aren't met.
type ServerMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ServerMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ServerMultiError) AllErrors() []error { return m }

// ServerValidationError is the validation error returned by Server.Validate if
// the designated constraints aren't met.
type ServerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ServerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ServerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ServerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ServerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ServerValidationError) ErrorName() string { return "ServerValidationError" }

// Error satisfies the builtin error interface
func (e ServerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sServer.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ServerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ServerValidationError{}

var _Server_Scheme_InLookup = map[string]struct{}{
	"http": {},
	"grpc": {},
	"tcp":  {},
}

// Validate checks the field values on Service with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Service) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Service with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ServiceMultiError, or nil if none found.
func (m *Service) ValidateAll() error {
	return m.validate(true)
}

func (m *Service) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !strings.HasPrefix(m.GetName(), "goim.") {
		err := ServiceValidationError{
			field:  "Name",
			reason: "value does not have prefix \"goim.\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !strings.HasSuffix(m.GetName(), ".service") {
		err := ServiceValidationError{
			field:  "Name",
			reason: "value does not have suffix \".service\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetVersion()) < 1 {
		err := ServiceValidationError{
			field:  "Version",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetLog()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Log",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Log",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLog()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "Log",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Metadata

	if all {
		switch v := interface{}(m.GetRedis()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Redis",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Redis",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRedis()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "Redis",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMq()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Mq",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Mq",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMq()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "Mq",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMysql()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Mysql",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ServiceValidationError{
					field:  "Mysql",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMysql()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "Mysql",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Http != nil {

		if all {
			switch v := interface{}(m.GetHttp()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ServiceValidationError{
						field:  "Http",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ServiceValidationError{
						field:  "Http",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetHttp()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ServiceValidationError{
					field:  "Http",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Grpc != nil {

		if all {
			switch v := interface{}(m.GetGrpc()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ServiceValidationError{
						field:  "Grpc",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ServiceValidationError{
						field:  "Grpc",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetGrpc()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ServiceValidationError{
					field:  "Grpc",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ServiceMultiError(errors)
	}

	return nil
}

// ServiceMultiError is an error wrapping multiple validation errors returned
// by Service.ValidateAll() if the designated constraints aren't met.
type ServiceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ServiceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ServiceMultiError) AllErrors() []error { return m }

// ServiceValidationError is the validation error returned by Service.Validate
// if the designated constraints aren't met.
type ServiceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ServiceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ServiceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ServiceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ServiceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ServiceValidationError) ErrorName() string { return "ServiceValidationError" }

// Error satisfies the builtin error interface
func (e ServiceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sService.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ServiceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ServiceValidationError{}

// Validate checks the field values on Log with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Log) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Log with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in LogMultiError, or nil if none found.
func (m *Log) ValidateAll() error {
	return m.validate(true)
}

func (m *Log) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if _, ok := Level_name[int32(m.GetLevel())]; !ok {
		err := LogValidationError{
			field:  "Level",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for EnableConsole

	if m.LogPath != nil {
		// no validation rules for LogPath
	}

	if len(errors) > 0 {
		return LogMultiError(errors)
	}

	return nil
}

// LogMultiError is an error wrapping multiple validation errors returned by
// Log.ValidateAll() if the designated constraints aren't met.
type LogMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LogMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LogMultiError) AllErrors() []error { return m }

// LogValidationError is the validation error returned by Log.Validate if the
// designated constraints aren't met.
type LogValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LogValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LogValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LogValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LogValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LogValidationError) ErrorName() string { return "LogValidationError" }

// Error satisfies the builtin error interface
func (e LogValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLog.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LogValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LogValidationError{}

// Validate checks the field values on Redis with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Redis) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Redis with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in RedisMultiError, or nil if none found.
func (m *Redis) ValidateAll() error {
	return m.validate(true)
}

func (m *Redis) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Addr

	// no validation rules for Password

	// no validation rules for MaxConns

	// no validation rules for MinIdleConns

	if d := m.GetDialTimeout(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = RedisValidationError{
				field:  "DialTimeout",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			lte := time.Duration(10*time.Second + 0*time.Nanosecond)
			gte := time.Duration(0*time.Second + 1000000*time.Nanosecond)

			if dur < gte || dur > lte {
				err := RedisValidationError{
					field:  "DialTimeout",
					reason: "value must be inside range [1ms, 10s]",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	if d := m.GetIdleTimeout(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = RedisValidationError{
				field:  "IdleTimeout",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			lte := time.Duration(10*time.Second + 0*time.Nanosecond)
			gte := time.Duration(0*time.Second + 1000000*time.Nanosecond)

			if dur < gte || dur > lte {
				err := RedisValidationError{
					field:  "IdleTimeout",
					reason: "value must be inside range [1ms, 10s]",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	if len(errors) > 0 {
		return RedisMultiError(errors)
	}

	return nil
}

// RedisMultiError is an error wrapping multiple validation errors returned by
// Redis.ValidateAll() if the designated constraints aren't met.
type RedisMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RedisMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RedisMultiError) AllErrors() []error { return m }

// RedisValidationError is the validation error returned by Redis.Validate if
// the designated constraints aren't met.
type RedisValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RedisValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RedisValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RedisValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RedisValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RedisValidationError) ErrorName() string { return "RedisValidationError" }

// Error satisfies the builtin error interface
func (e RedisValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRedis.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RedisValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RedisValidationError{}

// Validate checks the field values on MQ with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *MQ) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MQ with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MQMultiError, or nil if none found.
func (m *MQ) ValidateAll() error {
	return m.validate(true)
}

func (m *MQ) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for MaxRetry

	if len(errors) > 0 {
		return MQMultiError(errors)
	}

	return nil
}

// MQMultiError is an error wrapping multiple validation errors returned by
// MQ.ValidateAll() if the designated constraints aren't met.
type MQMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MQMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MQMultiError) AllErrors() []error { return m }

// MQValidationError is the validation error returned by MQ.Validate if the
// designated constraints aren't met.
type MQValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MQValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MQValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MQValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MQValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MQValidationError) ErrorName() string { return "MQValidationError" }

// Error satisfies the builtin error interface
func (e MQValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMQ.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MQValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MQValidationError{}

// Validate checks the field values on MySQL with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MySQL) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MySQL with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MySQLMultiError, or nil if none found.
func (m *MySQL) ValidateAll() error {
	return m.validate(true)
}

func (m *MySQL) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Addr

	// no validation rules for User

	// no validation rules for Password

	// no validation rules for Db

	// no validation rules for MaxIdleConns

	// no validation rules for MaxOpenConns

	if d := m.GetIdleTimeout(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = MySQLValidationError{
				field:  "IdleTimeout",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			lte := time.Duration(10*time.Second + 0*time.Nanosecond)
			gte := time.Duration(0*time.Second + 1000000*time.Nanosecond)

			if dur < gte || dur > lte {
				err := MySQLValidationError{
					field:  "IdleTimeout",
					reason: "value must be inside range [1ms, 10s]",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	if d := m.GetOpenTimeout(); d != nil {
		dur, err := d.AsDuration(), d.CheckValid()
		if err != nil {
			err = MySQLValidationError{
				field:  "OpenTimeout",
				reason: "value is not a valid duration",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		} else {

			lte := time.Duration(10*time.Second + 0*time.Nanosecond)
			gte := time.Duration(0*time.Second + 1000000*time.Nanosecond)

			if dur < gte || dur > lte {
				err := MySQLValidationError{
					field:  "OpenTimeout",
					reason: "value must be inside range [1ms, 10s]",
				}
				if !all {
					return err
				}
				errors = append(errors, err)
			}

		}
	}

	if len(errors) > 0 {
		return MySQLMultiError(errors)
	}

	return nil
}

// MySQLMultiError is an error wrapping multiple validation errors returned by
// MySQL.ValidateAll() if the designated constraints aren't met.
type MySQLMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MySQLMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MySQLMultiError) AllErrors() []error { return m }

// MySQLValidationError is the validation error returned by MySQL.Validate if
// the designated constraints aren't met.
type MySQLValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MySQLValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MySQLValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MySQLValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MySQLValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MySQLValidationError) ErrorName() string { return "MySQLValidationError" }

// Error satisfies the builtin error interface
func (e MySQLValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMySQL.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MySQLValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MySQLValidationError{}
