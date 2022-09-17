$ErrorActionPreference = 'Stop';

$packageName= 'kog'
$version = 'v0.0.1-alpha-04'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-386.zip"
$url64        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-amd64.zip"
$hash = '405f484ce30cf5363a75fc7273012234'
$hash64 = '483e95f90e9a9bccab28f00de9786e23'
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
