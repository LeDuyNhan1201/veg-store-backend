```shell
# Make IDE has full permission on current project
sudo chown -R $USER:$USER path/to/project
sudo chmod +x *.sh

docker rmi learn-go/go-clean-arch-builder
docker build -f docker/Dockerfile.builder -t learn-go/go-clean-arch-builder .

docker exec -it -uroot go-clean-arch bash
go run cmd/server/main.go
air -c .air.toml

docker exec -it -uroot go-clean-arch bash -c 'cd /app && air -c .air.toml'
docker exec -it -uroot go-clean-arch bash -c 'cd /app && go run cmd/server/main.go'

# Bonus
sudo ./start.sh
sudo ./stop.sh
```

```shell
go mod init <module-name>
go run <file-name>.go
go clean -modcache
go mod tidy
go get -u

go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n merge i18n/active.en.toml i18n/translate.vi.toml

go clean -modcache
git config --global http.sslVerify false
export GOINSECURE=github.com,go.googlesource.com,golang.org,go.uber.org,google.golang.org
export GOSUMDB=off
export GOPROXY=direct
GODEBUG=x509ignoreCN=1 go env -w GOINSECURE=*.org
go mod tidy
go get -u=patch
```
