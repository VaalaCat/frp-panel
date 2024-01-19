#!/usr/bin/env /bin/bash

OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64)
                wget -O frp-panel https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-amd64
                ;;
            aarch64)
                wget -O frp-panel https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-linux-arm64
                ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64)
                wget -O frp-panel https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-darwin-amd64
                ;;
            arm64)
                wget -O frp-panel https://github.com/VaalaCat/frp-panel/releases/latest/download/frp-panel-darwin-arm64
                ;;
        esac
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

chmod +x frp-panel

sudo mv frp-panel /usr/local/bin/frp-panel

get_start_params() {
    echo "请输入启动参数："
    read -r params
    echo "$params"
}

if [ -n "$1" ]; then
    start_params="$@"
else
    start_params=$(get_start_params)
fi

sudo tee /lib/systemd/system/frp-panel.service << EOF
[Unit]
Description=frp-panel
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/frp-panel $start_params

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload

sudo systemctl start frp-panel

sudo systemctl enable frp-panel
