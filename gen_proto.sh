#!/bin/bash
# define function
function gen_proto() {
    # scan all files in the giving directory
    # and generate proto files
    for file in $(ls $1)
    do
        # check if the file is a proto file
        if [[ $file == *.proto ]]
        then
            echo "generating $1/$file"
            # check if the file is a response.proto file
            if [[ $file == *response.proto ]]
            then
                # print the message in red color
                echo -e "\033[31m SPECIAL CASE $1/$file \033[0m"
                # generate response.proto file
                protoc --proto_path=. --proto_path=./third_party \
                		--go_out=paths=source_relative:. \
                		--go-grpc_out=paths=source_relative:. \
                		--validate_out=lang=go,paths=source_relative:. \
                		$1/$file
            else
                # generate other proto file
                protoc --proto_path=. --proto_path=./third_party --proto_path=./api \
                		--go_out=paths=source_relative:. \
                		--go-grpc_out=paths=source_relative:. \
                		--validate_out=lang=go,paths=source_relative:. \
                		$1/$file
            fi
        fi
        # check if the file is a directory and if it is not a hidden directory
        if [[ -d $1/$file && $file != .* ]]
        then
            # call the function recursively
            gen_proto $1/$file
        fi
    done
}

gen_proto $1