FROM docker.io/library/golang
RUN ["useradd", "--create-home", "user"]
USER user:user
RUN ["mkdir", "/home/user/page-change-mailer"]
WORKDIR /home/user/page-change-mailer
ENV PATH="$PATH":/go/bin
COPY --chown=user:user [".", "."]
RUN ["./docker.sh"]
