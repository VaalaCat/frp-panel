#!/usr/bin/env /bin/bash

OS=$(uname -s)
ARCH=$(uname -m)

if [ "$(curl -s https://ipinfo.io/country)" = "CN" ]; then
    prefix="https://mirror.ghproxy.com/"
    echo "监测到您的IP在中国，使用镜像下载"
else
    prefix=""
fi

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-amd64"
                ;;
            aarch64)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-arm64"
                ;;
            armv7l)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-armv7l"
                ;;
            armv6l)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-armv6l"
                ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-darwin-amd64"
                ;;
            arm64)
                wget -O frp-panel "${prefix}https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-darwin-arm64"
                ;;
        esac
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

chmod +x frp-panel

get_start_params() {
    read -p "请输入启动参数：" params
    echo "$params"
}

if [ -n "$1" ]; then
    start_params="$@"
else
    start_params=$(get_start_params)
fi

sudo ./frp-panel install $start_params

sudo systemctl daemon-reload

sudo systemctl start frpp

sudo systemctl restart frpp

sudo systemctl enable frpp