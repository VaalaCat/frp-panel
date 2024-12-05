# Download latest release from github
if($PSVersionTable.PSVersion.Major -lt 5){
    Write-Host "Require PS >= 5,your PSVersion:"$PSVersionTable.PSVersion.Major -BackgroundColor DarkGreen -ForegroundColor White
    exit
}
$clientrepo = "VaalaCat/frp-panel"
#  x86 or x64
if ([System.Environment]::Is64BitOperatingSystem) {
    if ([System.Environment]::Is64BitProcess) {
        $file = "frp-panel-windows-amd64.exe"
    } else {
        $file = "frp-panel-windows-arm64.exe"
    }
} else {
    Write-Host "Your system is 32-bit, please use 64-bit operating system" -BackgroundColor DarkGreen -ForegroundColor White
    exit
}
$clientreleases = "https://api.github.com/repos/$clientrepo/releases"
#重复运行自动更新
if (Test-Path "C:\frpp") {
    Write-Host "frp panel client already exists, delete and reinstall" -BackgroundColor DarkGreen -ForegroundColor White
    C:/frpp/frpp.exe stop
    C:/frpp/frpp.exe uninstall
    Start-Sleep -Seconds 3
    Remove-Item "C:\frpp\frpp.exe" -Recurse
}

#TLS/SSL
Write-Host "Determining latest frp panel client release" -BackgroundColor DarkGreen -ForegroundColor White
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$agenttag = (Invoke-WebRequest -Uri $clientreleases -UseBasicParsing | ConvertFrom-Json)[0].tag_name
#Region判断
$ipapi= Invoke-RestMethod  -Uri "https://api.myip.com/" -UserAgent "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1"
$region=$ipapi.cc
echo $ipapi
if($region -eq "CN" -or [string]::IsNullOrEmpty($ipapi)){
    $download = "https://ghp.ci/https://github.com/$clientrepo/releases/download/$agenttag/$file"
    Write-Host "Location:CN,use mirror address" -BackgroundColor DarkRed -ForegroundColor Green
}else{
    $download = "https://github.com/$clientrepo/releases/download/$agenttag/$file"
    Write-Host "Location:$region,connect directly!" -BackgroundColor DarkRed -ForegroundColor Green
}
echo $download
Invoke-WebRequest $download -OutFile "C:\frpp.exe"
New-Item -Path "C:\frpp" -ItemType Directory -ErrorAction SilentlyContinue
Move-Item -Path "C:\frpp.exe" -Destination "C:\frpp\frpp.exe"
C:\frpp\frpp.exe install $args
C:\frpp\frpp.exe start
Write-Host "Enjoy It!" -BackgroundColor DarkGreen -ForegroundColor Red