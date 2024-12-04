# 编译命令，根目录执行
KUBE_BUILD_PLATFORMS=linux/amd64 make all GOFLAGS=-v GOGCFLAGS="-N -l" -f build/root/Makefile