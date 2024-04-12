package host

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	IP_ADDR  = "IP_ADDR"
	HOSTNAME = "HOSTNAME"
)

var (
	HOSTNAME_REGEX = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
)

// HostInfo stores information retrieved for the user provided host
type HostInfo struct {
	IPs         []net.IP
	Hostnames   []string
	Nameservers []*net.NS
	Mailservers []*net.MX
}

func (i *HostInfo) Format() string {
	var formattedIps strings.Builder
	var formattedNameservers strings.Builder
	var formattedMailservers strings.Builder

	formattedHostnames := strings.Join(i.Hostnames, " ")

	for _, ip := range i.IPs {
		fmt.Fprintf(&formattedIps, "%s ", ip)
	}

	for _, ns := range i.Nameservers {
		fmt.Fprintf(&formattedNameservers, "%s ", ns.Host)
	}

	for _, mx := range i.Mailservers {
		fmt.Fprintf(&formattedMailservers, "%s ", mx.Host)
	}

	return fmt.Sprintf(
		"IPs: %s\nHostnames: %s\nNameservers: %s\nMailservers: %s",
		formattedIps.String(),
		formattedHostnames,
		formattedNameservers.String(),
		formattedMailservers.String(),
	)
}

// HostInfoCmd is used by [cmd/host.go] to get information
// about a given host.
type HostInfoCmd struct {
	Resolver *net.Resolver
}

func NewHostInfoCmd(ctx context.Context, cmd *cobra.Command) (*HostInfoCmd, error) {
	resolverAddress, err := cmd.Flags().GetString("resolver")

	if err != nil {
		return nil, err
	}

	// If there is no specific resolver passed as a command line
	// argument, use a default Resolver.
	if resolverAddress == "" {
		hostInfoCmd := &HostInfoCmd{
			Resolver: &net.Resolver{
				PreferGo: true,
			},
		}

		return hostInfoCmd, nil
	}

	// Ensure that the resolver is an IP address and not a hostname,
	// as the underlying Resolver.Dial method requires an IP address
	// to connect to.
	if ip := net.ParseIP(resolverAddress); ip == nil {
		return nil, errors.New(fmt.Sprintf("invalid resolver provided: %s. must be an IP address.", resolverAddress))
	}

	fmt.Println(fmt.Sprintf("using resolver: %s", resolverAddress))

	hostInfoCmd := &HostInfoCmd{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}

				// TODO: should the protocol and port also be configurable?
				return d.DialContext(ctx, "udp", fmt.Sprintf("%s:53", resolverAddress))
			},
		},
	}

	return hostInfoCmd, nil
}

func (h *HostInfoCmd) Run(ctx context.Context, host string) (HostInfo, error) {
	hostType, err := h.DetermineHostType(host)

	if err != nil {
		return HostInfo{}, err
	}

	var ips []net.IP
	var hostnames []string

	switch hostType {
	case IP_ADDR:
		ips = append(ips, net.ParseIP(host))
		hostnames, err = h.Resolver.LookupAddr(ctx, host)
	case HOSTNAME:
		hostnames = append(hostnames, host)
		ips, err = h.Resolver.LookupIP(ctx, "ip", host)
	default:
		return HostInfo{}, errors.New(fmt.Sprintf("cannot do lookup for unknown hostType %s", hostType))
	}

	if err != nil {
		return HostInfo{}, fmt.Errorf("error looking up host information for %s: %w", host, err)
	}

	nsRecords, err := h.getNSRecords(ctx, hostnames)

	if err != nil {
		return HostInfo{}, err
	}

	mxRecords, err := h.getMXReccords(ctx, hostnames)

	if err != nil {
		return HostInfo{}, err
	}

	hostInfo := HostInfo{
		IPs:         ips,
		Hostnames:   hostnames,
		Nameservers: nsRecords,
		Mailservers: mxRecords,
	}

	return hostInfo, nil
}

func (h *HostInfoCmd) getMXReccords(ctx context.Context, hostnames []string) ([]*net.MX, error) {
	var mxRecords []*net.MX

	for _, hostname := range hostnames {
		mailservers, err := h.Resolver.LookupMX(ctx, hostname)

		if err != nil {

			if dnsError, ok := err.(*net.DNSError); ok {

				// If the DNS error indicates that no records were found,
				// continue as the host might not have MX records configured.
				if dnsError.IsNotFound {
					continue
				}

			}

			// Otherwise, we have a real error and should return it
			return nil, fmt.Errorf("error looking up MX records for %s: %w", hostname, err)
		}

		mxRecords = append(mxRecords, mailservers...)
	}

	return mxRecords, nil
}

func (h *HostInfoCmd) getNSRecords(ctx context.Context, hostnames []string) ([]*net.NS, error) {
	var nsRecords []*net.NS

	for _, hostname := range hostnames {
		nameservers, err := h.Resolver.LookupNS(ctx, hostname)

		if dnsError, ok := err.(*net.DNSError); ok {

			// If the DNS error indicates that no records were found,
			// continue as the host might not have MX records configured.
			if dnsError.IsNotFound {
				continue
			}

			// Otherwise, we have a real error and should return it
			return nil, fmt.Errorf("error looking up NS records for %s: %w", hostname, err)
		}

		nsRecords = append(nsRecords, nameservers...)
	}

	return nsRecords, nil
}

func (h *HostInfoCmd) DetermineHostType(host string) (string, error) {
	if ip := net.ParseIP(host); ip != nil {
		return IP_ADDR, nil
	}

	if HOSTNAME_REGEX.MatchString(host) {
		return HOSTNAME, nil
	}

	return "", errors.New(fmt.Sprintf("could not determine host type of %s", host))
}
