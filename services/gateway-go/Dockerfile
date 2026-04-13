FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json .
RUN npm ci
COPY frontend .
RUN npm run build

FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod .
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /app/dist ./dist
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/dist ./dist
EXPOSE 8080
CMD ["./server"]
