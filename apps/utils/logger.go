package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type Context string

const RequestIDKey Context = "reqid"

type LoggerId struct {
	Sid string
	Ref string
}

func (l *LoggerId) SetSid(r *http.Request) {
	requestID, ok := r.Context().Value(RequestIDKey).(string)
	if !ok {
		KeyValue := ""
		session := uuid.NewV4()
		sessionSHA256 := session.String()
		KeyValue = strings.ReplaceAll(sessionSHA256, "-", "")
		l.Sid = KeyValue
	} else {
		l.Sid = requestID
	}
}

func (l *LoggerId) SetRef(ref any) {
	l.Ref = fmt.Sprintf("%v", ref)
}
func (l *LoggerId) RemoveRef() {
	l.Ref = ""
}

func (l *LoggerId) Log(message ...any) {
	log.Printf("%s - (%s) - %v\n", l.Sid, l.Ref, message)
}
