Write-Output $Env:GITHUB_REF
Write-Output $Env:GITHUB_ENV
Write-Output $Env.CHOCO_TOKEN

choco install windows-sdk-10-version-2104-all -y