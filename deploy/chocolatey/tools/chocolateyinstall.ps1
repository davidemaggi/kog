$ErrorActionPreference = 'Stop';

$packageName= 'kog'
$version = 'v0.0.1-alpha-02'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-386.zip"
$url        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-amd64.zip"
$hash = 'ebd1b3cdf4df932b66f5fe924c6978a0'
$hash64 = '025f9b74cbcc876c7e9ce0b22121f13e'
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
