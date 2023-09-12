# start
`docker compose -f "docker-compose.yaml" up -d --build `
`go run cmd/main.go`

# integrate

use the public code to start up the service as part of your api
use docker image

# passport

[![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![GoReport][report-img]][report] [![Coverage Status][cov-img]][cov]

Package `passport` is an identity access management api.

`go get github.com/georgi-georgiev/passport`

# database

create mongodb database `passport`

# docker
pull `docker pull bracer/passport`
run `docker run bracer/passport`

## Contributing

If you'd like to contribute to `passport`, we'd love your input! Please submit an issue first so we can discuss your proposal.

-------------------------------------------------------------------------------

Released under the [MIT License].

[MIT License]: LICENSE.txt
[doc-img]: https://pkg.go.dev/badge/github.com/georgi-georgiev/passport
[doc]: https://pkg.go.dev/github.com/georgi-georgiev/passport
[ci-img]: https://github.com/georgi-georgiev/passport/workflows/build/badge.svg
[ci]: https://github.com/georgi-georgiev/passport/actions
[report-img]: https://goreportcard.com/badge/github.com/georgi-georgiev/passport
[report]: https://goreportcard.com/report/github.com/georgi-georgiev/passport
[cov-img]: https://codecov.io/gh/georgi-georgiev/passport/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/georgi-georgiev/passport