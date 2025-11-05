## How to start application

# 1. Download & Setup all required packages
- Go - version: 1.25.3
- Docker (Builder, Engine)

# 2. Open terminal at root folder then run scripts (Must waiting for the first time)
```shell
sudo chown -R $USER:$USER .
sudo chmod +x scripts/*.sh
./scripts/start.sh
```

# 3. Run application in Docker container
```shell
docker exec -it -uroot go-clean-arch bash

# If you just want to run without Hot reload
go run cmd/server/main.go

# Or else, this way will help the app rebuild automatic when src code is updated
# Sometimes it not work so you must stop it using <Ctrl> <c> then run it again
air -c .air.toml
```

# 4. Test
Go to http://localhost:2345/api/v1/user/hello to make sure the app is working.
