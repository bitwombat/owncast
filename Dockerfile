FROM golang:1.19-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app .

FROM alpine:latest

RUN apk update && apk add curl wget python3
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod 755 /usr/local/bin/yt-dlp && \
    wget --progress=dot:mega https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz  && \
    tar xvf ffmpeg-release-amd64-static.tar.xz --strip-components=1 && \
    mv ffprobe ffmpeg /usr/local/bin && \
    rm ffmpeg-release-amd64-static.tar.xz

COPY --from=builder /usr/local/bin/app /usr/local/bin/app

# There's this for ffmpeg, but they're not static binaries
#    wget https://github.com/yt-dlp/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz && \
#
CMD ["/usr/local/bin/app"]
