Write-Host $Env:GITHUB_REF
Write-Host $Env:GITHUB_ENV
Write-Host $Env.CHOCO_TOKEN
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

choco install windows-sdk-10-version-2104-all -y