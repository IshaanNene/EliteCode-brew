#!/bin/bash

# Root files
touch go.mod go.sum main.go Dockerfile docker-compose.yml Makefile README.md .goreleaser.yml

# Formula
mkdir -p Formula
touch Formula/elitecode.rb

# cmd subdirectories and files
mkdir -p cmd/{auth,problems,user,system,github}
touch cmd/auth/{auth.go,login.go,logout.go,signup.go,whoami.go}
touch cmd/problems/{problems.go,list.go,set.go,run.go,submit.go,reset.go,search.go,bookmark.go,discuss.go}
touch cmd/user/{user.go,stats.go,my_problems.go}
touch cmd/system/{system.go,init.go,update.go}
touch cmd/github/{github.go,push.go}

# internal subdirectories and files
mkdir -p internal/{api,storage,docker/templates,github,utils}
touch internal/api/{client.go,auth.go,problems.go,user.go}
touch internal/storage/{config.go,cache.go}
touch internal/docker/runner.go
touch internal/docker/templates/{c.dockerfile,cpp.dockerfile,python.dockerfile,java.dockerfile,javascript.dockerfile}
touch internal/github/integration.go
touch internal/utils/{logger.go,file.go,spinner.go}

# backend subdirectories and files
mkdir -p backend/{routes,middleware,models,controllers,utils,docker}
touch backend/{package.json,package-lock.json,server.js,.env.example}
touch backend/routes/{auth.js,problems.js,users.js}
touch backend/middleware/{auth.js,validation.js}
touch backend/models/{User.js,Problem.js}
touch backend/controllers/{authController.js,problemController.js,userController.js}
touch backend/utils/{database.js,helpers.js}
touch backend/docker/{Dockerfile,docker-compose.yml}

# scripts
mkdir -p scripts
touch scripts/{build.sh,release.sh,install.sh}

echo "Directory structure created successfully."
