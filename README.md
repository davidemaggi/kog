# KOG
A KubeConfig manager for today's lazy developer

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

## Alias
In my day to day job I need to access many clusters, and they share very similar names(than you DevOps guys ❤) , let's give them some fancy new ones.

```shell
kog alias
```









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