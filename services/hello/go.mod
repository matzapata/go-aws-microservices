module github.com/matzapata/go-aws-microservices/services/hello

go 1.22.0

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/go-chi/chi/v5 v5.0.12
)

require shared/helpers v0.0.0-00010101000000-000000000000 // indirect

replace shared/helpers => ../../shared/helpers
