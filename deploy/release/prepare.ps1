Write-Output $Env:GITHUB_REF
Write-Output $Env:GITHUB_ENV
Write-Output $Env.CHOCO_TOKEN
Write-Output $Env.GITHUB_TOKEN



choco install windows-sdk-10-version-2104-all -y

choco install golang -y



dir "C:\Program Files (x86)\Windows Kits\10\bin\10.0. 19041.0\x64\"
dir "C:\Program Files (x86)\Windows Kits\10\bin\"