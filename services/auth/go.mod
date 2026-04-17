module mantevian.xyz/codenames/service_auth

go 1.26.2

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/lib/pq v1.12.3
	golang.org/x/crypto v0.50.0
	mantevian.xyz/codenames/shared v0.0.0
)

replace mantevian.xyz/codenames/shared => ../../shared
