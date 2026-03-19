package main

import "embed"

// webFS holds the compiled frontend assets from web/dist.
// The Dockerfile copies web/dist into cmd/api/web/dist before running go build
// so the embed path is valid and does not require a `..` traversal.
//
//go:embed web/dist
var webFS embed.FS
