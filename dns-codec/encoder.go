package dns

import "golang.org/x/net/dns/dnsmessage"

type RecordData interface {
	from() (dnsmessage.ResourceBody, error)
}

type Message struct {
	Header      Header
	Questions   []Questions
	Answers     []RR
	Authorities []RR
	Additionals []RR
}

type Header struct {
	ID                 uint16
	Response           bool
	Opcode             OpcodeType
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	AuthenticatedData  bool
	CheckingDisabled   bool
	RCode              RcodeType
}

type Questions struct {
	Name  string
	Type  Type
	Class ClassType
}

type RR struct {
	Name  string
	Type  Type
	Class ClassType
	TTL   uint32
	RData RecordData
}

func (m Message) Pack() ([]byte, error) {
	questions, err := m.toQuestions()
	if err != nil {
		return nil, err
	}

	answers, err := toResource(m.Answers)
	if err != nil {
		return nil, err
	}

	authorities, err := toResource(m.Authorities)
	if err != nil {
		return nil, err
	}

	additionals, err := toResource(m.Additionals)
	if err != nil {
		return nil, err
	}

	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:                 m.Header.ID,
			Response:           m.Header.Response,
			OpCode:             dnsmessage.OpCode(m.Header.Opcode),
			Authoritative:      m.Header.Authoritative,
			Truncated:          m.Header.Truncated,
			RecursionDesired:   m.Header.RecursionDesired,
			RecursionAvailable: m.Header.RecursionAvailable,
			AuthenticData:      m.Header.AuthenticatedData,
			CheckingDisabled:   m.Header.CheckingDisabled,
			RCode:              dnsmessage.RCode(m.Header.RCode),
		},
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additionals: additionals,
	}
	return msg.Pack()
}

func (m Message) toQuestions() ([]dnsmessage.Question, error) {
	data := make([]dnsmessage.Question, 0, len(m.Questions))
	for i := range m.Questions {
		name, err := newName(m.Questions[i].Name)
		if err != nil {
			return nil, err
		}

		data = append(data, dnsmessage.Question{
			Name:  name,
			Type:  dnsmessage.Type(m.Questions[i].Type),
			Class: dnsmessage.Class(m.Questions[i].Class),
		})
	}
	return data, nil
}

func toResource(rr []RR) ([]dnsmessage.Resource, error) {
	data := make([]dnsmessage.Resource, 0, len(rr))
	for i := range rr {
		name, err := newName(rr[i].Name)
		if err != nil {
			return nil, err
		}

		body, err := rr[i].RData.from()
		if err != nil {
			return nil, err
		}

		data = append(data, dnsmessage.Resource{
			Header: dnsmessage.ResourceHeader{
				Name:  name,
				Type:  dnsmessage.Type(rr[i].Type),
				Class: dnsmessage.Class(rr[i].Class),
				TTL:   rr[i].TTL,
			},
			Body: body,
		})
	}
	return data, nil
}

func newName(name string) (dnsmessage.Name, error) {
	return dnsmessage.NewName(name)
}
