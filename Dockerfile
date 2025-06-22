##############################
# 1 ── Build stage
##############################
FROM golang:1.24.2-alpine AS builder

# (Optional) bring in git & gcc if you compile CGO drivers
RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

# 1-a. cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# 1-b. copy the rest of the source
COPY . .
COPY .env .

# 1-c. build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


##############################
# 2 ── Runtime stage
##############################
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /src/server .
COPY --from=builder /src/.env   . 
################# ENVIRONMENT ###############
# Override these per-env (docker-compose, K8s, ECS ...)
ENV DB_HOST="host.docker.internal" \
    DB_USER="sagar" \
    DB_PORT=3306 \
    DB_PASS="sagar@30" \
    DB_NAME="video_analytics" \
    SERVER_PORT=3011 \
    REDIS_PORT="host.docker.internal:6379" \
    METRICS_PORT=9090

################# PORTS #####################
# main API
EXPOSE 3011   
# /metrics (Prometheus)
EXPOSE 9090  

################# ENTRYPOINT ###############
# Binary already listens on :8080 and spawns metrics on :9090
ENTRYPOINT ["/app/server"]
