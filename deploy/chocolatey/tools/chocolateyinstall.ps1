$ErrorActionPreference = 'Stop';

$packageName= 'kog'
$version = 'v0.0.1-alpha-02'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-amd64.zip"
$hash = '025f9b74cbcc876c7e9ce0b22121f13e'
$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = $url
  checksum      = $hash
  checksumType  = 'MD5'
}

Install-ChocolateyZipPackage @packageArgs
