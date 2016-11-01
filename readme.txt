# unofficial
docker build -t image_name .
docker create -v /go/images --name dataVol busybox /bin/true
docker run -d -p 8080:8080 --volumes-from=dataVol  --name container_name image_name

# ping test
localhost:8080/ping
