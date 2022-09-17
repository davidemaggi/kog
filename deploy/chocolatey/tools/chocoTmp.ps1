$ErrorActionPreference = 'Stop';

$packageName= 'kog'
$version = '@@VERSION@@'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-386-signed.zip"
$url64        = "https://github.com/davidemaggi/kog/releases/download/$version/kog-$version-windows-amd64-signed.zip"
$hash = '@@HASH_X86@@'
$hash64 = '@@HASH_X64@@'
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
