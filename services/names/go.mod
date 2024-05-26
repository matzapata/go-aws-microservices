module github.com/matzapata/go-aws-microservices/services/names

go 1.22.0

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.53.10
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/go-chi/chi/v5 v5.0.8
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	shared/helpers v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace shared/helpers => ../../shared/helpers
