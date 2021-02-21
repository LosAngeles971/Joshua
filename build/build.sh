if [[ -z $CI_PROJECT_DIR ]]
then
    echo "Missing CI_PROJECT_DIR"
    exit 1
fi
JOSHUA_WIN_64="joshua.exe"
JOSHUA_LINUX_64="joshua"
JOSHUA_LINUX_ARM="joshua.arm.i386"
JOSHUA_LINUX_ARM_64="joshua.arm.i64"

build_windows() {
    GOOS="windows" GOARCH="amd64" go build -o $CI_PROJECT_DIR/bin/$JOSHUA_WIN_64
}

build_linux() {
    GOOS="linux" GOARCH="amd64" go build -o $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_64
}

build_arm() {
    GOOS=linux GOARCH=arm GOARM=7 go build -o $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_ARM
    GOOS=linux GOARCH=arm64 go build -o $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_ARM_64
}

if [[ -z $1 ]]
then
    echo "Missing target"
    exit -1
fi

cd $CI_PROJECT_DIR/go
if [[ "$1" == "win" || "$1" == "all" ]]
then
    rm -f $CI_PROJECT_DIR/bin/$JOSHUA_WIN_64
    build_windows
fi
if [[ "$1" == "linux" || "$1" == "all" ]]
then
    rm -f $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_64
    build_linux
fi
if [[ "$1" == "arm" || "$1" == "all" ]]
then
    rm -f $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_ARM
    rm -f $CI_PROJECT_DIR/bin/$JOSHUA_LINUX_ARM_64
    build_arm
fi