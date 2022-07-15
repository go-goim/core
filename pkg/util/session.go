package util

import (
	"fmt"
	"strconv"

	messagev1 "github.com/go-goim/api/message/v1"
)

// Session generate session id
// switch tpye {
// case messagev1.SessionType_SingleChat:
// 	001fromUIDtoUID
// case messagev1.SessionType_GroupChat:
// 	010groupID00000000
// }
// to is groupID when type is groupChat
func Session(tp int32, from, to int64) string {
	// check if tp is valid
	if tp > 0xFF || tp < 0 {
		return ""
	}

	// int32 convert to uin8
	tp &= 0xFF
	u8 := uint8(tp)

	switch messagev1.SessionType(tp) {
	case messagev1.SessionType_SingleChat:
		if from > to {
			from, to = to, from
		}
		return fmt.Sprintf("%03d%020d%020d", u8, from, to)
	case messagev1.SessionType_GroupChat:
		return fmt.Sprintf("%03d%020d%020d", u8, to, 0)
	case messagev1.SessionType_Broadcast:
		return fmt.Sprintf("%03d%020d%020d", u8, 0, 0)
	case messagev1.SessionType_Channel:
		// channel session id is same as single chat
		// one of from and to is channel id
		if from > to {
			from, to = to, from
		}
		return fmt.Sprintf("%03d%020d%020d", u8, from, to)
	}

	return ""
}

var (
	// ErrInvalidSessionType invalid session type
	ErrInvalidSessionType = fmt.Errorf("invalid session type")
	// ErrInvalidSessionIDLength invalid session id length
	ErrInvalidSessionIDLength = fmt.Errorf("invalid session id length")
)

func ParseSession(s string) (tp int32, from, to int64, err error) {
	// check if s is valid
	if len(s) < 3+2*20 {
		return 0, 0, 0, ErrInvalidSessionIDLength
	}

	// first 3 bytes is session type
	tpStr := s[:3]
	i64, err := strconv.ParseInt(tpStr, 10, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	tp = int32(i64)

	// check if tp is valid
	if tp > 0xFF || tp < 0 {
		return 0, 0, 0, ErrInvalidSessionType
	}

	// get from and to
	fromStr := s[3 : 3+20]
	toStr := s[3+20 : 3+2*20]
	from, err = strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	to, err = strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	// if tp is group chat, to is group id
	if messagev1.SessionType(tp) == messagev1.SessionType_GroupChat {
		to = from
	}

	return tp, from, to, nil
}
