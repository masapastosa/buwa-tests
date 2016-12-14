package sockets

import (
	"crypto/rand"
	"errors"
)

type Session struct {
	ID int
	//SourceIP net.Addr
	Secret []byte
	// TODO: Add client struct to send the HTTP stream
}

type Sessions []Session

func NewSession(sessions Sessions) (*Session, error) {
	if cap(sessions) <= len(sessions) {
		return nil, errors.New("sessions: Sessions slice full")
	}

	ses := &Session{}
	ses.ID = len(sessions) + 1
	//ses.SourceIP = conn.RemoteAddr()
	ses.Secret = make([]byte, 10)
	if _, err := rand.Read(ses.Secret); err != nil {
		return nil, errors.New("session: Can't generate random value for session secret")
	}

	// Append to sessions slice
	sessions[ses.ID] = *ses
	return ses, nil
}
