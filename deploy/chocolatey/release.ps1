Set-Location .\deploy\chocolatey

Get-Location

choco install windows-sdk-10-version-2104-all -y
choco install gh -y
choco setapikey $env:CHOCO_TOKEN

#Get-ChildItem "C:\Program Files (x86)\Windows Kits\10\bin\"
$SignTool= "C:\Program Files (x86)\Windows Kits\10\bin\10.0.20348.0\x64\signtool.exe"



Get-ChildItem

$tag=$env:RELEASE
$tagStrip=$tag.substring(1)

$x64File="kog-$($tag)-windows-amd64.zip"
$x64Url= "https://github.com/davidemaggi/kog/releases/download/$($tag)/$($x64File)"

$x86File="kog-$($tag)-windows-386.zip"
$x86Url= "https://github.com/davidemaggi/kog/releases/download/$($tag)/$($x86File)"

# Downloaad the binaries
Invoke-WebRequest -URI $x64Url -OutFile $x64File
Invoke-WebRequest -URI $x86Url -OutFile $x86File

# Extract the binaries
$x86Dir=$x86File.replace('.zip','')
$x64Dir=$x64File.replace('.zip','')


Expand-Archive -LiteralPath $x64File -DestinationPath $x64Dir
Expand-Archive -LiteralPath $x86File -DestinationPath $x86Dir

Get-ChildItem

# Sign the Exe File

& $SignTool sign /f "cert.pfx" /p "$env:CODE_SIGN"  $x64Dir"\kog.exe"
& $SignTool sign /f "cert.pfx" /p "$env:CODE_SIGN" $x86Dir"\kog.exe"

Compress-Archive -Path $x64Dir\* -DestinationPath $x64Dir"-signed.zip"
Compress-Archive -Path $x86Dir\* -DestinationPath $x86Dir"-signed.zip"



#Create Choco Files
$nuspec=Get-Content -Path .\tmp.nuspec
$chocoScript=Get-Content -Path .\tools\chocoTmp.ps1

$nuspec=$nuspec.Replace("@@VERSION@@",$tagStrip)+""

Out-File -FilePath .\kog.nuspec -InputObject $nuspec


$md5x64= Get-FileHash $x64Dir"-signed.zip" -Algorithm MD5
$md5x86= Get-FileHash $x86Dir"-signed.zip" -Algorithm MD5

Out-File -FilePath $x64Dir"-signed.zip.md5" -InputObject $md5x64
Out-File -FilePath $x86Dir"-signed.zip.md5" -InputObject $md5x86

gh release upload $tag $x64Dir"-signed.zip"
gh release upload $tag $x64Dir"-signed.zip.md5"
gh release upload $tag $x86Dir"-signed.zip"
gh release upload $tag $x86Dir"-signed.zip.md5"

$chocoScript=$chocoScript.Replace("@@VERSION@@",$tag)+""
$chocoScript=$chocoScript.Replace("@@HASH_X64@@",$md5x64.Hash)+""
$chocoScript=$chocoScript.Replace("@@HASH_X86@@",$md5x86.Hash)+""

Out-File -FilePath .\tools\chocolateyinstall.ps1 -InputObject $chocoScript
Get-ChildItem

choco pack kog.nuspec

choco push "kog."$tagStrip".nupkg"

