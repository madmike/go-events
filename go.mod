module github.com/madmike/go-events

go 1.25.5

require (
	github.com/google/uuid v1.6.0
	github.com/madmike/go-infra v0.0.1
	github.com/nats-io/nats.go v1.51.0
)

require (
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rs/zerolog v1.35.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
)

replace github.com/madmike/go-infra => ../infra
