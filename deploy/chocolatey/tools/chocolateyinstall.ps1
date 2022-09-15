$ErrorActionPreference = 'Stop';

$packageName= 'kog'
$version = 'v0.0.1-alpha-04'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-386.zip"
$url64        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-amd64.zip"
$hash = '9c858767d268faf03d7051e263d8cb04'
$hash64 = 'cf6551402e9d06e8635fba1afc40bb48'
$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url64bit      = $url64
  url           = $url
  checksum      = $hash
  checksum64    = $hash64
  checksumType  = 'MD5'
  checksumType64  = 'MD5'
}

Install-ChocolateyZipPackage @packageArgs
