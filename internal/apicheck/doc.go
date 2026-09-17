// Package apicheck compares cloud.google.com/go/firestore against fsmock.
//
// Layers:
//   - discover: every SDK type with exported methods must be in the registry or ignoredTypes
//   - signature: each SDK method must exist on the fsmock interface with a normalized
//     signature (or an explicit deviations/exceptions entry with a Reason)
//   - coverage: when FSMOCK_COVERPROFILE is set, every *xxxWrapper method must be
//     covered by tests or listed in waivers
//
// Stale exceptions/deviations/ignoredTypes/waivers fail the build.
package apicheck
