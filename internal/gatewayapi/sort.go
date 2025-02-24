// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package gatewayapi

import (
	"sort"

	"github.com/envoyproxy/gateway/internal/ir"
)

type XdsIRRoutes []*ir.HTTPRoute

func (x XdsIRRoutes) Len() int      { return len(x) }
func (x XdsIRRoutes) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x XdsIRRoutes) Less(i, j int) bool {
	// 1. Sort based on characters in a matching path.
	pCountI := pathMatchCount(x[i].PathMatch)
	pCountJ := pathMatchCount(x[j].PathMatch)
	if pCountI < pCountJ {
		return true
	}
	if pCountI > pCountJ {
		return false
	}
	// Equal case

	// 2. Sort based on the number of Header matches.
	hCountI := len(x[i].HeaderMatches)
	hCountJ := len(x[j].HeaderMatches)
	if hCountI < hCountJ {
		return true
	}
	if hCountI > hCountJ {
		return false
	}
	// Equal case

	// 3. Sort based on the number of Query param matches.
	qCountI := len(x[i].QueryParamMatches)
	qCountJ := len(x[j].QueryParamMatches)
	return qCountI < qCountJ
}

// sortXdsIR sorts the xdsIR based on the match precedence
// defined in the Gateway API spec.
// https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1beta1.HTTPRouteRule
func sortXdsIRMap(xdsIR XdsIRMap) {
	for _, ir := range xdsIR {
		for _, http := range ir.HTTP {
			// descending order
			sort.Sort(sort.Reverse(XdsIRRoutes(http.Routes)))
		}
	}
}

func pathMatchCount(pathMatch *ir.StringMatch) int {
	if pathMatch != nil {
		if pathMatch.Exact != nil {
			return len(*pathMatch.Exact)
		}
		if pathMatch.Prefix != nil {
			return len(*pathMatch.Prefix)
		}
		if pathMatch.SafeRegex != nil {
			return len(*pathMatch.SafeRegex)
		}
	}
	return 0
}
