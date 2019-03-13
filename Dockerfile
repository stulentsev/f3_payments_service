#######################
# Stage: builder
FROM golang:1.11-alpine AS builder
LABEL description="Base build image used by other stages"
LABEL maintainer="Nikolai Vladimirov <nikolay@vladimiroff.com>"


# Install some dependencies needed to build the project
RUN apk add bash git gcc g++ libc-dev make
WORKDIR /home/app/services/f3_payments_service

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Cache modules
RUN go mod download

#######################
# Stage: runner
# This image holds app code for general usage
FROM builder AS runner
RUN go mod download
# Here we copy the rest of the source code
# Make sure your .dockerignore file is properly configure to ensure proper layer caching
COPY . .

#######################
# Stage: server_builder
# This image is used to build the server executable
FROM runner AS server_builder

# And compile the project
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/f3paymentsd

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS production
# Finally we copy the statically compiled Go binary.
COPY --from=server_builder /go/bin/f3paymentsd /bin/f3paymentsd

EXPOSE 8080
ENTRYPOINT ["/bin/f3paymentsd"]
