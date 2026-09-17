package fsmock

// Versioning scheme: v2.<firestore-minor>.<patch>
//
// The second number of this module's semver always matches the minor version of
// cloud.google.com/go/firestore that the release was tested against.
// Example: fsmock v2.25.x pairs with firestore v1.25.x.
//
// Patch versions are free for fsmock-only fixes. There is no backport policy —
// only the latest SupportedFirestoreMinor receives fixes. Older tags remain
// immutable on the Go module proxy for consumers pinned to older SDKs.
//
// Caveat: Go modules have no upper bound. The firestore version in go.mod is a
// floor (minimum), not a pin. MVS may select a newer firestore in a consumer
// module. "Paired with" means "tested against", not "exclusive lock".
// See COMPATIBILITY.md.

// Version is the fsmock release version string (no leading "v").
const Version = "2.25.0"

// FirestoreSDKVersion is the exact cloud.google.com/go/firestore version this
// release was developed and tested against.
const FirestoreSDKVersion = "v1.25.0"

// SupportedFirestoreMinor is the Firestore SDK minor paired with this release.
// Must equal the minor in FirestoreSDKVersion and the minor of Version.
const SupportedFirestoreMinor = 25
