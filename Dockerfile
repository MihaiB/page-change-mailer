FROM docker.io/library/golang:alpine AS build
RUN ["adduser", "--disabled-password", "user"]
RUN ["apk", "update"]
RUN ["apk", "upgrade"]
RUN ["apk", "add", "bash", "gcc", "libc-dev"]
USER user:user
RUN ["mkdir", "/home/user/page-change-mailer"]
WORKDIR /home/user/page-change-mailer
COPY --chown=user:user [".", "."]
RUN ["./docker.sh"]

FROM docker.io/library/alpine
RUN ["adduser", "--disabled-password", "user"]
USER user:user
COPY --from=build --chown=user:user ["/go/bin/page-change-mailer", "/home/user"]
RUN ["mkdir", "/home/user/data"]
CMD ["/home/user/page-change-mailer"]
