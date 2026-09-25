package dns

import (
	"fmt"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type AResource struct {
	A net.IP
}

func (r *AResource) from() (dnsmessage.ResourceBody, error) {
	if ip4 := net.IP.To4(r.A); ip4 != nil {
		var arr [4]byte
		copy(arr[:], ip4)
		return &dnsmessage.AResource{A: arr}, nil
	}
	return nil, fmt.Errorf("not a valid IPv4 address: %v", r.A)
}

type AAAAResource struct {
	AAAA net.IP
}

func (r *AAAAResource) from() (dnsmessage.ResourceBody, error) {
	if ip16 := net.IP.To16(r.AAAA); ip16 != nil {
		var arr [16]byte
		copy(arr[:], ip16)
		return &dnsmessage.AAAAResource{AAAA: arr}, nil
	}
	return nil, fmt.Errorf("not a valid IPv6 address: %v", r.AAAA)
}

type MXResource struct {
	Pref uint16
	MX   string
}

func (r *MXResource) from() (dnsmessage.ResourceBody, error) {
	name, err := dnsmessage.NewName(r.MX)
	if err != nil {
		return nil, err
	}

	return &dnsmessage.MXResource{
		Pref: r.Pref,
		MX:   name,
	}, nil
}

type CNAMEResource struct {
	CNAME string
}

func (r *CNAMEResource) from() (dnsmessage.ResourceBody, error) {
	name, err := dnsmessage.NewName(r.CNAME)
	if err != nil {
		return nil, err
	}

	return &dnsmessage.CNAMEResource{CNAME: name}, nil
}

type PTRResource struct {
	PTR string
}

func (r *PTRResource) from() (dnsmessage.ResourceBody, error) {
	name, err := dnsmessage.NewName(r.PTR)
	if err != nil {
		return nil, err
	}

	return &dnsmessage.PTRResource{PTR: name}, nil
}

type NSResource struct {
	NS string
}

func (r *NSResource) from() (dnsmessage.ResourceBody, error) {
	name, err := dnsmessage.NewName(r.NS)
	if err != nil {
		return nil, err
	}

	return &dnsmessage.NSResource{NS: name}, nil
}
