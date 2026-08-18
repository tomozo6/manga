package main

import "database/sql"

// app contains the dependencies used to serve authenticated manga requests.
type app struct {
	verifier       tokenVerifier
	allowed        map[string]struct{}
	gcsMediaSigner gcsMediaSigner
	db             *sql.DB
}
