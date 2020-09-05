FROM docker.io/library/golang AS build
RUN ["useradd", "--create-home", "user"]
USER user:user
RUN ["mkdir", "/home/user/page-change-mailer"]
WORKDIR /home/user/page-change-mailer
COPY --chown=user:user [".", "."]
RUN ["./docker.sh"]

FROM docker.io/library/golang
RUN ["useradd", "--create-home", "user"]
USER user:user
COPY --from=build --chown=user:user ["/go/bin/page-change-mailer", "/go/bin"]
RUN ["mkdir", "/home/user/data"]
ENV PATH="$PATH":/go/bin
CMD ["page-change-mailer"]
