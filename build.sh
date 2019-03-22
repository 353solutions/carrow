docker build . -t carrow:builder
docker run -v $PWD:/home/carrow -it --workdir=/home/carrow/ carrow:builder