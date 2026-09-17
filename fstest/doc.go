// Package fstest provides Firestore emulator helpers and integration/parity tests.
//
// Correctness layers:
//
//   - apicheck: SDK surface completeness + signature parity (CI unit job)
//
//   - unit: unwrap / nil / panic-recovery / readSettings aliasing contracts
//
//   - accuracy: scripts/check-accuracy.sh — integration tests with coverage;
//     any fstest skip fails the gate; every *xxxWrapper method must be covered
//     or listed in internal/apicheck/waivers.go
//
//     make emulator-up
//     make accuracy
package fstest
