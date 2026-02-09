package sessionId

import (
	"strconv"
	"strings"
)

type SessionId string

func NewSessionId(name, token string, id uint64) SessionId {
	if len(name) == 0 || len(token) == 0 {
		panic("name or token are required")
	}
	idstr := strconv.FormatUint(id, 10)
	return SessionId(name + ":" + token + ":" + idstr)
}

func FromString(str string) SessionId {
	return SessionId(str)
}

func (s SessionId) Name() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	return arr[0]
}

func (s SessionId) Token() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	return arr[1]
}

func (s SessionId) Id() uint64 {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}
	return id
}

func (s SessionId) Info() (string, string, uint64) {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}
	return arr[0], arr[1], id
}

func (s SessionId) String() string {
	return string(s)
}
