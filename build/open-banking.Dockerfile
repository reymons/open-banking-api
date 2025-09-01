FROM golang:alpine AS base
WORKDIR /app
RUN apk update && apk add make

FROM base AS build
COPY . .
RUN make build-open-banking

FROM base AS final
COPY --from=build /app/open-banking.bin .
CMD ["./open-banking.bin"]

