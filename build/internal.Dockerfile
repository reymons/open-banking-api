FROM golang:alpine AS base
WORKDIR /app
RUN apk update && apk add make

FROM base AS build
COPY . .
RUN make build-internal

FROM base AS final
COPY --from=build /app/bin/internal .
CMD ["./internal"]
