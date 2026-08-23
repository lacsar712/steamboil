FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/steamboil ./cmd/steamboil

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/steamboil /steamboil
EXPOSE 8080
ENTRYPOINT ["/steamboil"]
