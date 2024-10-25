package v1

import (
	"encoding/json"
	"errors"
)

const version = "1"

type MessageType int

const (
	MessageNewChallenge MessageType = iota
	MessageSolvedChallenge
	MessageWordOfWisdom
)

type Message struct {
	ProtocolVersion string
	Type            MessageType
	TypedMessage    any
}

type NewChallengeMessage struct {
	Challenge  string
	Difficulty int
}

func NewNewChallengeMessage(challenge string, difficulty int) Message {
	return Message{
		ProtocolVersion: version,
		Type:            MessageNewChallenge,
		TypedMessage: NewChallengeMessage{
			Challenge:  challenge,
			Difficulty: difficulty,
		},
	}
}

type SolvedChallengeMessage struct {
	Challenge string
	Solution  string
}

func NewSolvedChallengeMessage(challenge string, solution string) Message {
	return Message{
		ProtocolVersion: version,
		Type:            MessageSolvedChallenge,
		TypedMessage: SolvedChallengeMessage{
			Challenge: challenge,
			Solution:  solution,
		},
	}
}

type WordOfWisdomMessage struct {
	WordOfWisdom string
}

func NewWordOfWisdomMessage(quote string) Message {
	return Message{
		ProtocolVersion: version,
		Type:            MessageWordOfWisdom,
		TypedMessage: WordOfWisdomMessage{
			WordOfWisdom: quote,
		},
	}
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var raw struct {
		ProtocolVersion string
		Type            MessageType
		TypedMessage    json.RawMessage
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.ProtocolVersion = raw.ProtocolVersion
	m.Type = raw.Type

	switch m.Type {
	case MessageNewChallenge:
		var msg NewChallengeMessage
		if err := json.Unmarshal(raw.TypedMessage, &msg); err != nil {
			return err
		}
		m.TypedMessage = msg

	case MessageSolvedChallenge:
		var msg SolvedChallengeMessage
		if err := json.Unmarshal(raw.TypedMessage, &msg); err != nil {
			return err
		}
		m.TypedMessage = msg

	case MessageWordOfWisdom:
		var msg WordOfWisdomMessage
		if err := json.Unmarshal(raw.TypedMessage, &msg); err != nil {
			return err
		}
		m.TypedMessage = msg

	default:
		return errors.New("unknown message type")
	}

	return nil
}
