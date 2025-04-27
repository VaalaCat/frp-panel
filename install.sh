#!/usr/bin/env /bin/bash

OS=$(uname -s)
ARCH=$(uname -m)

if ping -c 1 -W 1 google.com > /dev/null 2>&1; then
    prefix=""
    echo "检测到您的网络可以连接到 Google，不使用镜像下载"
else
    prefix="https://ghfast.top/"
    echo "检测到您的网络无法连接到 Google，使用镜像下载"
fi

current_dir=$(pwd)
temp_dir=$(mktemp -d)
echo "下载临时文件夹创建在: $temp_dir"
cd "$temp_dir"

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

sudo chmod +x frp-panel

cd "$current_dir"

new_executable_path="$temp_dir/frp-panel"

get_start_params() {
    read -p "请输入启动参数：" params
    echo "$params"
}

find_frpp_executable() {
    service_file=$(systemctl show -p FragmentPath frpp.service 2>/dev/null | cut -d'=' -f2)
    if [[ -z "$service_file" || ! -f "$service_file" ]]; then
        echo ""
        return 1
    fi
    exec_start=$(grep -oP '^ExecStart=\K.*' "$service_file")
    if [[ -z "$exec_start" ]]; then
        echo ""
        return 1
    fi
    executable_path=$(echo "$exec_start" | awk '{print $1}')
    echo "$executable_path"
}

if [ -n "$1" ]; then
    start_params="$@"
else
    start_params=$(get_start_params)
fi

if systemctl list-units --type=service | grep -q frpp; then
    echo "frpp 服务存在"
    executable_path=$(find_frpp_executable)
    if [ -z "$executable_path" ]; then
        echo "无法找到 frpp 服务的执行文件路径，请检查systemd文件"
        exit 1
    fi
    echo "更新程序到原路径：$executable_path"
    sudo rm -rf "$executable_path"
    sudo cp "$new_executable_path" "$executable_path"
    $executable_path version
    sudo $executable_path uninstall
    sudo $executable_path install $start_params
    sudo systemctl daemon-reload
    echo "参数已重写，请执行 cat /etc/systemd/system/frpp.service 仔细检查启动命令，避免无法启动"
    echo "执行 sudo systemctl restart frpp 重启服务"
    echo "执行 sudo systemctl status frpp 查看服务状态"
    exit 0
else
    echo "frpp 服务不存在，进行安装"
fi

sudo cp "$new_executable_path" .

sudo ./frp-panel install $start_params

echo "frp-panel 服务安装完成, 安装路径：$(pwd)/frp-panel"

sudo systemctl daemon-reload

sudo ./frp-panel start

sudo ./frp-panel version

echo "frp-panel 服务已启动"

sudo systemctl restart frpp

sudo systemctl enable frpp
