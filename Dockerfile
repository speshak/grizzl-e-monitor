FROM golang:1.26-alpine AS build

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

RUN apk add --no-cache make git && \
    make build

# Create runtime image
FROM alpine:3.24
COPY --from=build /app/build/grizzl-e-monitor /bin/grizzl-e-monitor

EXPOSE 8080
CMD ["/bin/grizzl-e-monitor"]
