git pull 
go build main.go
mv -f main /opt/logservice2/server
supervisorctl restart logservice2-api
