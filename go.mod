module github.com/josephburnett/jd

go 1.24.0

toolchain go1.24.13

require (
	github.com/go-openapi/jsonpointer v0.22.4
	github.com/josephburnett/jd/v2 v2.0.0-20240818191833-6125a15c637a
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/go-openapi/swag/jsonname v0.25.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/josephburnett/jd/v2 => ./v2
