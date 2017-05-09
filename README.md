# match
 a match server build on [kitfw](https://github.com/yangwenhai/kitfw)

## zipkin:
wget -O zipkin.jar 'https://search.maven.org/remote_content?g=io.zipkin.java&a=zipkin-server&v=LATEST&c=exec'
java -jar zipkin.jar   

# build  

$ cd match  && source devenv.sh

$ cd src/kitfw/vendor && govendor sync

$ cd ../../../

$ go build -o match kitfw/match/server


# run

./run.sh

