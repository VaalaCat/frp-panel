#!/usr/bin/env /bin/bash

OS=$(uname -s)
ARCH=$(uname -m)

# --- Argument Parsing ---
custom_proxy=""
frp_panel_args=() # Use an array to hold arguments for frp-panel

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --github-proxy)
            if [ -z "$2" ]; then
                echo "Error: --github-proxy requires a URL argument."
                exit 1
            fi
            custom_proxy="$2"
            echo "使用自定义的 GitHub 镜像: $custom_proxy"
            shift 2 # shift past argument name and value
            ;;
        *)
            # Collect remaining arguments for frp-panel
            frp_panel_args+=("$1") # Add the argument to the array
            shift 1 # shift past argument
            ;;
    esac
done
# --- End Argument Parsing ---


# --- Network Check and Prefix Setting ---
prefix=""
if ping -c 1 -W 1 google.com > /dev/null 2>&1; then
    echo "检测到您的网络可以连接到 Google，不使用镜像下载"
else
    if [ -n "$custom_proxy" ]; then
        prefix="$custom_proxy"
        echo "检测到您的网络无法连接到 Google，使用用户提供的镜像下载: $prefix"
    else
        prefix="https://ghfast.top/"
        echo "检测到您的网络无法连接到 Google，使用默认镜像下载: $prefix"
    fi
fi
# --- End Network Check and Prefix Setting ---


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
            *)
                 echo "Unsupported Linux architecture: $ARCH"
                 exit 1
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
            *)
                echo "Unsupported Darwin architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

if [ ! -f frp-panel ]; then
    echo "Error: Download failed. frp-panel executable not found."
    exit 1
fi

sudo chmod +x frp-panel

cd "$current_dir"

new_executable_path="$temp_dir/frp-panel"

# Function to find the frpp executable path from systemd service file
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

# --- Install or Update frp-panel ---
if systemctl list-units --type=service | grep -q frpp; then
    echo "frpp 服务存在，进行更新"
    executable_path=$(find_frpp_executable)
    if [ -z "$executable_path" ]; then
        echo "无法找到 frpp 服务的执行文件路径，请检查systemd文件"
        exit 1
    fi
    echo "更新程序到原路径：$executable_path"
    sudo rm -f "$executable_path" # Use -f to avoid prompt if file is read-only
    sudo cp "$new_executable_path" "$executable_path"
    sudo "$executable_path" version
    sudo "$executable_path" uninstall
    # Pass collected frp_panel_args to install command
    sudo "$executable_path" install "${frp_panel_args[@]}"
    sudo systemctl daemon-reload
    echo "参数已重写，请执行 cat /etc/systemd/system/frpp.service 仔细检查启动命令，避免无法启动"
    echo "执行 sudo systemctl restart frpp 重启服务"
    echo "执行 sudo systemctl status frpp 查看服务状态"
    exit 0
else
    echo "frpp 服务不存在，进行安装"
fi

# Copy the new executable to the current directory for initial install
sudo cp "$new_executable_path" .

# Run the install command with collected frp_panel_args
sudo ./frp-panel install "${frp_panel_args[@]}"

echo "frp-panel 服务安装完成, 安装路径：$(pwd)/frp-panel"

sudo systemctl daemon-reload

sudo ./frp-panel start

sudo ./frp-panel version

echo "frp-panel 服务已启动"

sudo systemctl restart frpp

sudo systemctl enable frpp

# Clean up temporary directory
rm -rf "$temp_dir"
echo "清理临时文件夹: $temp_dir"
