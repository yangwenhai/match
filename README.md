# match
 a match server build on [kitfw](https://github.com/yangwenhai/kitfw)

## zipkin:
wget -O zipkin.jar 'https://search.maven.org/remote_content?g=io.zipkin.java&a=zipkin-server&v=LATEST&c=exec'
java -jar zipkin.jar   

## development

1、define a struct in message.go  and run compile.sh
2、add message protocol id in protocol.go
3、add protocol id to HandlerMap in match/service/service_define.go
4、define a Handler（like SumHandler、NewSumHandler） in match/service/service_define.go and implement Process
5、add a new handler file （like match/service/sum_handler.go ）and implement doProcess

## bambam and  capnpc-go
1、capnp
wget https://github.com/sandstorm-io/capnproto/archive/v0.6.0.tar.gz
tar xvf v0.6.0.tar.gz 
cd capnproto-0.6.0/c++/cmake
cmake ../
make
ll src/capnp/

2、bambam
go get github.com/shurcooL/go-goon
git clone https://github.com/glycerine/bambam.git 
cd bambam && make

3、capnpc-go（ 'capnp compile -ogo ./schema.capnp' will need capnpc-go）
go get github.com/glycerine/go-capnproto 
git clone https://github.com/glycerine/go-capnproto 
cd go-capnproto && make

# build  

$ cd match  && source devenv.sh

$ cd src/kitfw/vendor && govendor sync

$ cd ../../../

$ go build -o match kitfw/match/server


# run

./run.sh

