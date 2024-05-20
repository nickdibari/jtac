package subdomains

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"strings"
	"sync"
)

// TODO: Allow for user-supplied wordlist as command line option?

//go:embed subdomain_wordlist.txt
var wordlist string

// BruteForceResult is used by `checkCandidateSubdomain` and SubdomainInfoCmd.bruteForce
// to hold information the results of a brute force search. It needs a mutex in order to
// avoid race conditions when keeping track of valid subdomains for a given host.
type BruteForceResult struct {
	Mutex  sync.Mutex
	Result []string
}

// SubdomainInfoCmd is used by [cmd/subdomain.go] to enumerate
// possible subdomains for a given host.
type SubdomainInfoCmd struct{}

func (s *SubdomainInfoCmd) Run(ctx context.Context, host string) ([]string, error) {
	var result []string

	bruteForcedSubdomains, err := s.bruteForce(ctx, host)

	if err != nil {
		return []string{}, fmt.Errorf("error doing brute force subdomain search for %s: %w", host, err)
	}

	result = append(result, bruteForcedSubdomains...)

	// TODO: Find way to dedupe the return value?
	return result, nil
}

func (s *SubdomainInfoCmd) bruteForce(ctx context.Context, host string) ([]string, error) {
	var waitGroup sync.WaitGroup

	candidateSubdomains := strings.Split(wordlist, "\n")

	result := BruteForceResult{}

	for _, subdomain := range candidateSubdomains {
		candidate := fmt.Sprintf("%s.%s", subdomain, host)

		waitGroup.Add(1)

		go checkCandidateSubdomain(&waitGroup, candidate, &result)
	}

	waitGroup.Wait()

	return result.Result, nil
}

func checkCandidateSubdomain(waitGroup *sync.WaitGroup, candidate string, result *BruteForceResult) {
	defer waitGroup.Done()

	// If we didn't get an error for a DNS lookup for the candidate
	// domain (subdomain from wordlist + given host), then we can
	// assume that the candidate is a valid subdomain for the given
	// host.
	//
	// TODO: Find way to confirm that this subdomain is in fact a
	// valid subdomain for the given host. Does it point to something
	// that actually takes traffic? Does it just redirect to the main
	// domain?
	if _, err := net.LookupHost(candidate); err == nil {
		result.Mutex.Lock()
		defer result.Mutex.Unlock()

		result.Result = append(result.Result, candidate)
	}
}
