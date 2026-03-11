# boot4go-DockerUI
A visual management tools for docker container and docker swarm cluster, You can browse and maintain the docker single node 
or cluster node both worker and Manager.

![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

## INTRO
*DockerUI* is an easy-to-use and lightweight docker management tool. Through the operation of the web interface, it is more 
convenient for users unfamiliar with docker instructions to enter the docker world faster.

*DockerUI* has an easy-to-use interface. It does not need to remember docker instructions. Just download the image 
And you can immediately join and complete the deployment. Based on the features of docker, the version of the image can 
be directly updated in docker UI. With the same settings, the upgrade can be completed by redeploying and replacing the 
riginal container, and the functions of the latest version can be used.

*DockerUI* covers more than 95% of the command functions of the docker cli command line. Through the visual operation function 
provided in the dockerui interface, the management and maintenance functions of the docker environment and the docker swarm 
cluster environment can be easily performed.

*DockerUI* is a visual graphical management tool for docker container images. Dockerui can be used to easily build, manage 
and maintain the docker environment. It is completely open-source and free. It is based on the container installation method, 
which is convenient and efficient for deployment.

Our primary goals are:
You can use Docker UI to manage the docker and swarm instance more easy.

Official site: https://github.com/jonnyan404/docker-ui


## Feature

- *Docker host managementmanagement*
  Data volume management, image management, container management, build management, warehouse configuration management, network configuration management

- *Docker swarm cluster management*
  Cluster profile information, node management, service management, task management, password management, configuration management

- *Task arrangement*
  Docker task scheduling, docker swarm task scheduling

## Snapshot

### Home page (summary)
![image](https://img-blog.csdnimg.cn/46f6144e1f2b4562bc6c55b98d19b018.png)

### Image list
![image](https://img-blog.csdnimg.cn/c2161006901e4ee7bee132dc1271bc44.png)

### Search repository / pull image
![image](https://img-blog.csdnimg.cn/97519b029c8d4910bd21401e271f2b71.png)

### Build Image
![image](https://img-blog.csdnimg.cn/b23fdb3c295e4ecdbaa3904cc79e7c8b.png)

### Export / Import Image
![image](https://img-blog.csdnimg.cn/ebcf17ec4203495bafe6e599da9891fc.png)

### Push Image
![image](https://img-blog.csdnimg.cn/c8e55811bc234ed9893cc8f9f1ba4a5a.png)

#### Execute Image
![image](https://img-blog.csdnimg.cn/4c25eaeaa7b14d07838cf10c04ead5fd.png)

### List Container
![image](https://img-blog.csdnimg.cn/c6fe99139d654eed885748d5c86070b1.png)

### Web Console of Container
![image](https://img-blog.csdnimg.cn/28ec5d0ce20945db9908153d7090ac77.png)

### Container File System
![image](https://img-blog.csdnimg.cn/b464d6d70e534b4087ab67cd4a381f27.png)

### Stats of Container
![image](https://img-blog.csdnimg.cn/f34cef2f67e442b09e1cdf31a6907f07.png)

### List processes of Container
![image](https://img-blog.csdnimg.cn/a4204e7673294ed1bb01a03798beb823.png)

### Export file from Container
![image](https://img-blog.csdnimg.cn/a4204e7673294ed1bb01a03798beb823.png)

### Export file from Container
![image](https://images.gitee.com/uploads/images/2022/0530/104343_11e0da56_6575697.png)

### Network Management
![image](https://img-blog.csdnimg.cn/09f53d750e054911876cf7f1b44da520.png)

### Swarm Cluster Management
![image](https://img-blog.csdnimg.cn/b8ee779df1e141968042be7a77c7bbf6.png)

### Create Service
![image](https://img-blog.csdnimg.cn/842d0e8f5b3f4f3c968c5b7ca099cd8d.png)

### Task Management
![image](https://img-blog.csdnimg.cn/ac521683f92a4b1098fb87d93d66c134.png)

### List Task
![image](https://img-blog.csdnimg.cn/d7684158e28c42eb830d33180fa86be4.png)

### Docker Compose
![image](https://img-blog.csdnimg.cn/45e6e9185a9a4ac4888f90130547d93f.png)


## Installation and Getting Started

### From Github
- Download sourcecode from github website, visit https://github.com/jonnyan404/docker-ui .
- Install the golang runtime environment.
- Come into the project directory
- Run command as blow;
  - export GO111MODULE=on 
  - export GOPROXY="https://goproxy.cn,direct"
  - go mod tidy
  - go mod download
  - go build -o docker-ui .
- Run ./docker-ui command to start

### From GitHub Release (Binaries)

If you don't want to install Go locally, you can download pre-built binaries from GitHub Releases.

1) Download the binary for your OS/CPU from the Release assets, for example:
- Linux: `docker-ui-linux-amd64`, `docker-ui-linux-arm64`, `docker-ui-linux-armv7`
- macOS: `docker-ui-darwin-amd64`, `docker-ui-darwin-arm64`
- Windows: `docker-ui-windows-amd64.exe`, `docker-ui-windows-arm64.exe`

2) Release binaries embed the static UI files, so they can run as a single self-contained executable.

3) Run it:
- Linux/macOS:
  - `chmod +x ./docker-ui-linux-amd64`
  - `./docker-ui-linux-amd64 --addr :8999 --endpoint unix`
- Windows (PowerShell):
  - `./docker-ui-windows-amd64.exe --addr :8999 --endpoint unix`

Tip: You may rename the binary to `docker-ui` for convenience.

### Form docker-compose

```
# docker-compose.yml
services:
    docker-ui:
        container_name: docker-ui
        restart: always
        volumes:
            - ./data:/app/config
            - /var/run/docker.sock:/var/run/docker.sock
        ports:
            - 8999:8999
        image: jonnyan404/docker-ui
        #command: -l 0.0.0.0:9000
```

### From docker
- pull image from hub
  - docker image pull jonnyan404/docker-ui
- start container with image, and publish 8999 port to your port
  - docker run -d --name docker-ui -v /var/run/docker.sock:/var/run/docker.sock -p 8999:8999 jonnyan404/docker-ui

## Visit the browser tool
- Now, you can visit like as http://192.168.56.102:8999 .
- Default Username/Password dockerui/dockerui
- Enjoy it now.

## Command line options

All available command line options:

| Option | Short | Default | Description |
|---|---:|---:|---|
| `--addr` | `-l` | `:8999` | Listen address of DockerUI HTTP server. |
| `--issuer` | `-i` | `DEFAULT_ISSUER` | JWT token issuer. |
| `--token_expire` | `-e` | `24` | Token expiration time in hours. |
| `--endpoint` |  | `unix` | Docker endpoint. Use `unix` for local `/var/run/docker.sock`, or `HOST:PORT` for TCP (e.g. `192.168.56.102:2375`). |
| `--license` |  | empty | CubeUI license string. |
| `--reset-user` |  | empty | Reset password for username (creates user if not exists), then exit unless `--reset-keep-running` is set. |
| `--reset-password` |  | empty | Password used with `--reset-user`. |
| `--reset-keep-running` |  | `false` | Keep running the server after reset/create user. |
| `--help` | `-h` |  | Show help. |
| `--version` |  |  | Show version. |

## Forgot password / bootstrap admin

If you forgot the password and cannot login to the UI, you can reset (or create) a user in the local sqlite database (`./data.db`) before starting the server.

- Reset password (or create if not exists), then exit:
  - `go run . --reset-user dockerui --reset-password NEW_PASSWORD`
- Reset password (or create if not exists) using binary, then exit (rename your binary to `dockerui` first):
- Reset password (or create if not exists) using binary, then exit:
  - `./docker-ui --reset-user dockerui --reset-password NEW_PASSWORD`
- Keep running the server after reset/create:
  - `go run . --reset-user dockerui --reset-password NEW_PASSWORD --reset-keep-running`
- Keep running the server after reset/create using binary:
  - `./docker-ui --reset-user dockerui --reset-password NEW_PASSWORD --reset-keep-running`

Note:
- `dockerui` is created only on first database initialization (fresh `./data.db`). If you delete it later, it will NOT be restored automatically on restart.
- Use `--reset-user` / `--reset-password` to recover access when no user can login.