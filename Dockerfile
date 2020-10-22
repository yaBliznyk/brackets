FROM golang:1.15-alpine

ENV PROJECT_NAME brackets
ENV SOURCE_PATH /go/src/github.com/yaBliznyk/brackets
ENV BUILD_PATH /etc

RUN mkdir -pv $BUILD_PATH && mkdir -pv $SOURCE_PATH

ADD ./ $SOURCE_PATH

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o $BUILD_PATH/$PROJECT_NAME $SOURCE_PATH/cmd/$PROJECT_NAME/main.go

RUN rm -rf $SOURCE_PATH
RUN chmod +x $BUILD_PATH/$PROJECT_NAME

EXPOSE 8080

CMD ["/etc/brackets"]