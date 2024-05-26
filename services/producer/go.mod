module github.com/matzapata/go-aws-microservices/services/producer

go 1.22.0

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.53.10
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/go-chi/chi/v5 v5.0.12
	shared/helpers v0.0.0-00010101000000-000000000000
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

replace shared/helpers => ../../shared/helpers
