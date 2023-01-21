# go-vendors-api
Go GraphQL (and REST) API for storing vendors

## Updating Graphql Schemas
The `github.com/99designs/gqlgen` package was used to auto generate files for graphql. The schemas live in `graph/schema.graphqls`. Use `go generate` (after removing schema.resolvers.go and navigating to the graph dir) to generate new `schema.resolvers.go` and `model/models_gen.go` files.
