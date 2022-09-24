# KOG
A KubeConfig manager for the lazy developer of today

## Wizard
The main kog command, it will prompt the available Contexts and Namespaces to let you switch between them in a lazy way...
oh, and it allows you to filter... in case you have a gazillion of contexts
```shell
kog
```
### Vanilla
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_wizard1.gif?raw=true)
### With Filter
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_wizard2.gif?raw=true)

## Merge
>With great projects comes many clusters (cit. uncle Ben).

This Command will merge the desired yaml file in your kubeconfig file
```shell
kog merge your/file/path/config.yaml
```
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_merge.gif?raw=true)
## Alias
In my day to day job I need to access many clusters, and they share very similar names(than you DevOps guys ❤) , let's give them some fancy new ones.

```shell
kog alias
```
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_alias.gif?raw=true)

## Info
Display the current Context and namespace.
Or the raw yaml file if you are brave enough
```shell
kog info
```
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_info1.gif?raw=true)
```shell
kog info --raw
```
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_info2.gif?raw=true)


## Version
Display the current kog version
```shell
kog --version
```

![](https://github.com/davidemaggi/kog/blob/0dd7bc010eccf465cfbef1c01a2ddf79254b1086/imgs/readme/gif_version.gif?raw=true)

## Help
Display the help view, if you are too lazy to read the previous instructions...
In such case, we are the same... I like you...
```shell
kog help
```
![](https://github.com/davidemaggi/kog/blob/main/imgs/readme/gif_help.gif?raw=true)

## Global Flags

| Flag      | Alias | Optonal | Description                                                                    | Default               |
|-----------|:-----:|:-------:|--------------------------------------------------------------------------------|-----------------------|
| --config  |  -c   |    x    | The config file to manage| $USERDIR/.kube/config |
| --force   |  -f   |    x    | Perform the operation forcing it even if no actual action is needed            | false                 |
| --verbose |  -v   |    x    | Display extended logs                                                          | false                 |

## Install

Kog is making its first baby steps, It's available as standalone binaries for all major platforms that you can download and install by yourself, the only package manager from where you can download kog is Chocolatey on windows.
>Why? because Ii mainly work on windows that's why
 
>⭐: Preferred Method
### macOS 🍎
#### Manual 🔨
1. Download the binary that fits your OS from [Release](https://github.com/davidemaggi/kog/releases) page
2. Copy kog executable to an accessible directory on your Mac(e.g /usr/local/bin)
3. Add the Folder to your PATH Environment variable, if not already configured
4. Enjoy Kog

### Linux 🐧
#### Manual 🔨
1. Download the binary that fits your OS from [Release](https://github.com/davidemaggi/kog/releases) page
2. Copy kog.exe file to an accessible directory on your PC(e.g /usr/bin)
3. Add the Folder to your PATH Environment variable, if not already configured
4. Enjoy Kog
5. 
### Windows 🪟
#### Manual 🔨
1. Download the binary that fits your OS from [Release](https://github.com/davidemaggi/kog/releases) page
2. Copy kog.exe file to an accessible directory on your PC(e.g C:\Users\<YourUserName>\AppData\Local)
3. Add the Folder to your PATH Environment variable, if not already configured
4. Enjoy Kog
#### Chocolatey 📦 ⭐
Simply execute the choco command to install it
```shell
choco install kog
```
To update it
```shell
choco upgrade kog
```