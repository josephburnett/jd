module github.com/josephburnett/jd

go 1.18

require (
	github.com/go-openapi/jsonpointer v0.19.5
	github.com/josephburnett/jd/v2 v2.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
)

replace github.com/josephburnett/jd/v2 => ./v2
