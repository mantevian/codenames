module mantevian.xyz/codenames/service_gateway

go 1.26.2

require mantevian.xyz/codenames/shared v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

replace mantevian.xyz/codenames/shared => ../../../shared
