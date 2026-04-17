`sudo docker compose --file compose.dev.yaml up --build`

# env
```
POSTGRES_URL=postgres://<user>:<password>@postgres:5432/codenames?sslmode=disable
POSTGRES_DB=codenames
POSTGRES_USER=<user>
POSTGRES_PASSWORD=<password>
PGADMIN_DEFAULT_EMAIL=<email>
PGADMIN_DEFAULT_PASSWORD=<admin_password>
JWT_SECRET=<jwt_secret>
```