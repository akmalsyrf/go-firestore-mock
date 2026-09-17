// Package fstest provides Firestore emulator helpers and integration/parity tests.
//
// Correctness layers:
//
//   - apicheck: SDK surface completeness (CI unit job)
//
//   - unit: unwrap / nil / panic-recovery contracts
//
//   - accuracy: scripts/check-accuracy.sh — required parity + integration;
//     t.Skip on a required test fails the gate
//
//     export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080
//     make accuracy
package fstest
