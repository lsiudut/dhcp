package dhcpv6

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/insomniacslk/dhcp/rfc1035label"
)

// WithClientID adds a client ID option to a DHCPv6 packet
func WithClientID(duid Duid) Modifier {
	return func(d DHCPv6) DHCPv6 {
		cid := OptClientId{Cid: duid}
		d.UpdateOption(&cid)
		return d
	}
}

// WithServerID adds a client ID option to a DHCPv6 packet
func WithServerID(duid Duid) Modifier {
	return func(d DHCPv6) DHCPv6 {
		sid := OptServerId{Sid: duid}
		d.UpdateOption(&sid)
		return d
	}
}

// WithNetboot adds bootfile URL and bootfile param options to a DHCPv6 packet.
func WithNetboot(d DHCPv6) DHCPv6 {
	msg, ok := d.(*DHCPv6Message)
	if !ok {
		log.Printf("WithNetboot: not a DHCPv6Message")
		return d
	}
	// add OptionBootfileURL and OptionBootfileParam
	opt := msg.GetOneOption(OptionORO)
	if opt == nil {
		opt = &OptRequestedOption{}
	}
	// TODO only add options if they are not there already
	oro := opt.(*OptRequestedOption)
	oro.AddRequestedOption(OptionBootfileURL)
	oro.AddRequestedOption(OptionBootfileParam)
	msg.UpdateOption(oro)
	return d
}

// WithUserClass adds a user class option to the packet
func WithUserClass(uc []byte) Modifier {
	// TODO let the user specify multiple user classes
	return func(d DHCPv6) DHCPv6 {
		ouc := OptUserClass{UserClasses: [][]byte{uc}}
		d.AddOption(&ouc)
		return d
	}
}

// WithArchType adds an arch type option to the packet
func WithArchType(at iana.Arch) Modifier {
	return func(d DHCPv6) DHCPv6 {
		ao := OptClientArchType{ArchTypes: []iana.Arch{at}}
		d.AddOption(&ao)
		return d
	}
}

// WithIANA adds or updates an OptIANA option with the provided IAAddress
// options
func WithIANA(addrs ...OptIAAddress) Modifier {
	return func(d DHCPv6) DHCPv6 {
		opt := d.GetOneOption(OptionIANA)
		if opt == nil {
			opt = &OptIANA{}
		}
		iaNa := opt.(*OptIANA)
		for _, addr := range addrs {
			iaNa.AddOption(&addr)
		}
		d.UpdateOption(iaNa)
		return d
	}
}

// WithDNS adds or updates an OptDNSRecursiveNameServer
func WithDNS(dnses ...net.IP) Modifier {
	return func(d DHCPv6) DHCPv6 {
		odns := OptDNSRecursiveNameServer{
			NameServers: append([]net.IP{}, dnses[:]...),
		}
		d.UpdateOption(&odns)
		return d
	}
}

// WithDomainSearchList adds or updates an OptDomainSearchList
func WithDomainSearchList(searchlist ...string) Modifier {
	return func(d DHCPv6) DHCPv6 {
		osl := OptDomainSearchList{
			DomainSearchList: &rfc1035label.Labels{
				Labels: searchlist,
			},
		}
		d.UpdateOption(&osl)
		return d
	}
}

// WithRequestedOptions adds requested options to the packet
func WithRequestedOptions(optionCodes ...OptionCode) Modifier {
	return func(d DHCPv6) DHCPv6 {
		opt := d.GetOneOption(OptionORO)
		if opt == nil {
			opt = &OptRequestedOption{}
		}
		oro := opt.(*OptRequestedOption)
		for _, optionCode := range optionCodes {
			oro.AddRequestedOption(optionCode)
		}
		d.UpdateOption(oro)
		return d
	}
}
