package util

import (
	"fmt"
	"strconv"
	"strings"

	messagev1 "github.com/go-goim/api/message/v1"
	"github.com/go-goim/core/pkg/types"
)

// Session generate session id
// switch tpye {
// case messagev1.SessionType_SingleChat:
// 	001fromUIDtoUID
// case messagev1.SessionType_GroupChat:
// 	010groupID00000000
// }
// to is groupID when type is groupChat
func Session(tp messagev1.SessionType, from, to types.ID) string {
	// check if tp is valid
	if tp > 0xFF || tp < 0 {
		return ""
	}

	// int32 convert to uin8
	tp &= 0xFF
	u8 := strconv.FormatInt(int64(tp), 16)

	switch tp {
	case messagev1.SessionType_SingleChat:
		if from > to {
			from, to = to, from
		}
		return fmt.Sprintf("%02s%011s%011s", u8, from.Base58(), to.Base58())
	case messagev1.SessionType_GroupChat:
		return fmt.Sprintf("%02s%011s%011d", u8, from.Base58(), 0)
	case messagev1.SessionType_Broadcast:
		return fmt.Sprintf("%02s%011d%011d", u8, 0, 0)
	case messagev1.SessionType_Channel:
		// channel session id is same as single chat
		// one of from and to is channel id
		if from > to {
			from, to = to, from
		}
		return fmt.Sprintf("%02s%011s%011s", u8, from.Base58(), to.Base58())
	default:
		// in default case, use single chat rule
		// only return empty when tp is invalid
		if from > to {
			from, to = to, from
		}
		return fmt.Sprintf("%02s%011s%011s", u8, from.Base58(), to.Base58())
	}
}

var (
	// ErrInvalidSessionType invalid session type
	ErrInvalidSessionType = fmt.Errorf("invalid session type")
	// ErrInvalidSessionIDLength invalid session id length
	ErrInvalidSessionIDLength = fmt.Errorf("invalid session id length")
)

func ParseSession(s string) (tp messagev1.SessionType, from, to types.ID, err error) {
	// check if s is valid
	if len(s) < 2+2*11 {
		return 0, 0, 0, ErrInvalidSessionIDLength
	}

	// first 2 bytes is session type
	tpStr := s[:2]
	i32, err := strconv.ParseInt(tpStr, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}
	tp = messagev1.SessionType(i32)

	// check if tp is valid
	if tp > 0xFF || tp < 0 {
		return 0, 0, 0, ErrInvalidSessionType
	}

	// get from and to
	// trim 0 from string from and to
	// in some case, the first char is 0, so we need to trim it
	fromStr := strings.TrimLeft(s[2:2+11], "0")
	toStr := strings.TrimLeft(s[2+11:2+2*11], "0")
	from, err = types.ParseBase58([]byte(fromStr))
	if err != nil {
		return 0, 0, 0, err
	}
	to, err = types.ParseBase58([]byte(toStr))
	if err != nil {
		return 0, 0, 0, err
	}

	// if tp is group chat, to is group id
	if tp == messagev1.SessionType_GroupChat {
		to = from
	}

	return tp, from, to, nil
}
