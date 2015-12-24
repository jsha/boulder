// Copyright 2015 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bdns

import (
	"fmt"
	"net"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/miekg/dns"
	"github.com/letsencrypt/boulder/probs"
)

type dnsError struct {
	rCode      int
	recordType uint16
	hostname   string
	// TODO(jsha@eff.org): The authority section is available in the Ns field of
	// dns.Msg. We should include that in this struct and print it as part of the
	// error detail.
}

func (d dnsError) Error() string {
	return fmt.Sprintf("DNS problem: %s looking up %s for %s", dns.RcodeToString[d.rCode],
		dns.TypeToString[d.recordType], d.hostname)
}

const detailDNSTimeout = "DNS query timed out"
const detailDNSNetFailure = "DNS networking error"
const detailServerFailure = "Server failure at resolver"

// ProblemDetailsFromDNSError checks the error returned from Lookup...
// methods and tests if the error was an underlying net.OpError or an error
// caused by resolver returning SERVFAIL or other invalid Rcodes and returns
// the relevant core.ProblemDetails.
func ProblemDetailsFromDNSError(err error) *probs.ProblemDetails {
	problem := &probs.ProblemDetails{Type: probs.ConnectionProblem}
	if netErr, ok := err.(*net.OpError); ok {
		if netErr.Timeout() {
			problem.Detail = detailDNSTimeout
		} else {
			problem.Detail = detailDNSNetFailure
		}
	} else if dnsErr, ok := err.(*dnsError); ok {
		problem.Detail = dnsErr.Error()
	} else {
		problem.Detail = detailServerFailure
	}
	return problem
}
