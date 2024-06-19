protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/*/*.proto
protoc --go_out=../ --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=import api/ptypes/src/*.proto
#protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=import api/*/src/*.proto
protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=import api/games/*/src/*.proto

for file in api/*
do
    if test -d $file
    then
        dir_name=$(basename $file)
        if [ $dir_name != "ptypes" -a $dir_name != "errors"  -a $dir_name != "games" ]; then
               protoc  -I . --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=import -I ./api/thirdparty  --grpc-gateway_out=. --grpc-gateway_opt=paths=import  api/${dir_name}/src/*.proto
        fi
    fi
done
