Write-Output $Env:GITHUB_REF
Write-Output $Env:GITHUB_ENV
Write-Output $Env.CHOCO_TOKEN
Write-Output $Env.GITHUB_TOKEN



choco install windows-sdk-10-version-2104-all -y

choco install golang -y


Get-ChildItem "C:\Program Files (x86)\Windows Kits\10\bin\"
Get-ChildItem "C:\Program Files (x86)\Windows Kits\10\bin\10.0. 2104.0\x64\"
