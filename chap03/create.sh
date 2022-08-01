docker rmi httpserver
docker build -t  httpserver:1.0 .
docker rm -f httpserver
docker run -itd -p 8010:80 -p 8360:8360 --name httpserver --restart=always httpserver:1.0 
#docker push httpserver:1.0
PID=$(docker inspect --format "{{ .State.Pid }}" httpserver)
nsenter -t $PID -n ip a
curl http://10.1.0.34:8010/healthz
