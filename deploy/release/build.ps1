# Build Windows x64

$Env:GOOS = "windows";
$Env:GOARCH = "amd64";

go build -o kog.exe


dir

Write-Output $Env:GITHUB_REF
Write-Output $Env:GITHUB_ENV
Write-Output $Env.CHOCO_TOKEN
Write-Output $Env.GITHUB_TOKEN