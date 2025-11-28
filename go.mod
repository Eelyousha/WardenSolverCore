module wardensolvercore

go 1.21

require (
	github.com/lib/pq v1.10.9
	handlerv1 v0.0.0-00010101000000-000000000000
)

replace (
	generators => ./generators
	handlerv1 => ./handlers/v1
)
