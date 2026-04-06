(
  seq 10000 | xargs -n1 -P100 -I{} curl -s -o /dev/null http://localhost:8080/api/v1/health/
) &

(
  seq 10000 | xargs -n1 -P100 -I{} curl -s -o /dev/null http://localhost:8080/api/v1/health/ping
) &

wait