composerize docker run -v /home/csh0101/lab/netagent/optl/alloy.config:/etc/alloy/config.alloy -p 12345:12345 grafana/alloy:latest run --stability.level=public-preview --server.http.listen-addr=0.0.0.0:12345 --storage.path=/var/lib/alloy/data /etc/alloy/config.alloy

composerize docker run -p 4040:4040 grafana/pyroscope
