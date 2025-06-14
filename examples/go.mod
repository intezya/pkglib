module examples

go 1.24.2

replace (
	github.com/intezya/pkglib/configloader => ./../configloader
	github.com/intezya/pkglib/crypto => ./../crypto
	github.com/intezya/pkglib/generate => ./../generate
	github.com/intezya/pkglib/itertools => ./../itertools
	github.com/intezya/pkglib/logger => ./../logger
	github.com/intezya/pkglib/sliceutils => ./../sliceutils
	github.com/intezya/pkglib/jwtlib => ./../jwtlib
)

require (
	github.com/intezya/pkglib/configloader v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/crypto v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/generate v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/itertools v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/logger v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/sliceutils v0.0.0-00010101000000-000000000000
	github.com/intezya/pkglib/jwtlib v0.0.0-00010101000000-000000000000
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1 // indirect
)
