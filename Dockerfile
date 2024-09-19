FROM golang:1.23-alpine AS build
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .
 
RUN make build
 

# Create runtime image
FROM alpine:3.20
COPY --from=build /app/build/grizzl-e-prom /bin/grizzl-e-prom

EXPOSE 8080
CMD ["/bin/grizzl-e-prom"]