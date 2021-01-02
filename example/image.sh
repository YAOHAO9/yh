rm -rf vendor # 删除之前的依赖包

go get -u # 更新依赖

go mod vendor # 将依赖放在项目内，统一打包到docker镜像,后面就可以不用安装直接使用了

docker build . -t pine # build镜像

docker tag pine 192.168.200.193:5000/pine # 打标签

docker push 192.168.200.193:5000/pine # 存放到私有仓库