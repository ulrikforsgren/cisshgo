module github.com/tbotnz/cisshgo

go 1.17

require (
	github.com/gliderlabs/ssh v0.3.5
	golang.org/x/crypto v0.5.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/term v0.4.0 // indirect
)

replace golang.org/x/term v0.4.0 => ./src/term // indirect
