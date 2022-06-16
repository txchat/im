# golang1.17 or latest
TARGETDIR=target

projectVersion=$(shell git describe --abbrev=8 --tags)
gitCommit=$(shell git rev-parse --short=8 HEAD)

pkgCommitName=${projectVersion}_${gitCommit}
servers=comet logic

help: ## 查看makefile帮助文档
	@printf "Help doc:\nUsage: make [command]\n"
	@printf "[command]\n"
	@grep -h -E '^([a-zA-Z_-]|\%)+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: clean ## 编译本机系统和指令集的可执行文件
	./script/build/builder.sh ${TARGETDIR}

build_%: clean ## 编译目标机器的可执行文件（例如: make build_linux_amd64）
	./script/build/builder.sh ${TARGETDIR} $*

pkg: build ## 编译并打包本机系统和指令集的可执行文件
	tar -zcvf ${TARGETDIR}_'host'_${pkgCommitName}.tar.gz ${TARGETDIR}/

pkg_%: build_% ## 编译并打包目标机器的可执行文件（例如: make pkg_linux_amd64）
	tar -zcvf ${TARGETDIR}_$*_${pkgCommitName}.tar.gz ${TARGETDIR}/

images: build_linux_amd64 ## 打包docker镜像
	cp script/docker/*Dockerfile ${TARGETDIR}
	cd ${TARGETDIR} && for i in $(servers) ; do \
		docker build . -f $$i.Dockerfile -t txchat-$$i:${projectVersion}; \
	done

init-compose: images ## 使用docker compose启动
	cp -R script/compose/. run_compose/
	cd run_compose && \
	./initwork.sh "${servers}" "${projectVersion}"

docker-compose-up:  ## 使用docker compose启动
	@if [ ! -d "run_compose/" ]; then \
		exit -1;\
	 fi; \
	cd run_compose && \
	docker compose -f components.compose.yaml -f service.compose.yaml up -d

docker-compose-%: ## 使用docker compose 命令(服务列表：make docker-compose-ls；停止服务：make docker-compose-stop；卸载服务：make docker-compose-down)
	@if [ ! -d "run_compose/" ]; then \
       cp -R script/compose/. run_compose/; \
     fi; \
    cd run_compose && \
    docker compose -f components.compose.yaml -f service.compose.yaml $*

test:
	$(GOENV) go test -v ./...

clean:
	rm -rf ${TARGETDIR}

run:
	for i in $(servers) ; do \
		nohup ${TARGETDIR}/$$i -conf=${TARGETDIR}/$$i.toml -logtostderr 2>&1 > ${TARGETDIR}/$$i.log; \
	done

stop:
	for i in $(servers) ; do \
	  	pkill -f ${TARGETDIR}/$$i; \
	done
