FROM  alpine:3.7

ENV APP_HOME /home/exporter

ADD exporter*.tar.gz ${APP_HOME}

RUN chmod a+x ${APP_HOME}/*.sh

RUN chmod 755 ${APP_HOME}/*.sh

EXPOSE 8080

WORKDIR ${APP_HOME}

CMD ${APP_HOME}/run.sh
