FROM alpine
ADD song-service /song-service
ENTRYPOINT [ "/song-service" ]
