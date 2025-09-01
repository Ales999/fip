module github.com/ales999/fip

go 1.24.6

//replace github.com/ales999/cisaccs/v2 => ../cisaccs/v2

require (
	github.com/alecthomas/kong v1.12.1
	github.com/ales999/cisaccs/v2 v2.0.0-20250901033336-574e529ba168
)

require (
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
