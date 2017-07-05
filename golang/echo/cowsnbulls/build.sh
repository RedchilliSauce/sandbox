mkdir -p data
cd utils
./getalljsons.sh
cd ..
go get github.com/labstack/echo
go get github.com/RedchilliSauce/sandbox/sandbox/golang/echo/cowsnbulls/utils

go build -o bin/application application.go