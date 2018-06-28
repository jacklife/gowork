FROM zenap/alpine-glibc:3.4

ENV APP_HOME /home/exporter

ADD exporter*.tar.gz ${APP_HOME}

RUN chmod 755 ${APP_HOME}/*.sh

EXPOSE 8080

WORKDIR ${APP_HOME}

CMD ${APP_HOME}/run.sh
