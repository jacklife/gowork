Prometheus exporter for PostgreSQL server metrics. Supported PostgreSQL versions: 9.0 and up to 10th.

Build and run

mvn clean install -P  windows(linux)

./exporter.exe

#docker

docker build -t exporter:1.0 .

docker run -ti  -p 9292:9292 -d  exporter:1.0 /bin/sh