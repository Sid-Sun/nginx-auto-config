# nginx-auto-config
## GoLang program to configure nginx virtual servers interactively 

### Presets available:

1: Create static site config (with index).

2: Create config to host files (w/o index)

3: Create config for Angular/Vue production site with routing

4: Proxy pass requests to a port or a site

5: Serve a PHP site with fastcgi and php-fpm

6: Permanent URL redirection to someplace else.

7: Configure to forward all HTTP requests to HTTPS

8: Port forward without hostname with custom port numbers

### Compiled binaries:

> [Linux amd64 / x86_64](https://cdn.sidsun.com/nginx-auto-config/nginx-auto-config_linux-amd64)

### Debian Packages:

> [amd64](https://cdn.sidsun.com/nginx-auto-config/nginx-auto-config.deb)

### Use YAPPA ( Yet Another Personal Package Archive ):

```bash
curl -s --compressed "https://sid-sun.github.io/yappa/KEY.gpg" | sudo apt-key add -
curl -s --compressed "https://sid-sun.github.io/yappa/yappa.list" | sudo tee /etc/apt/sources.list.d/yappa.list
sudo apt update
sudo apt install nginx-auto-config
```

:)