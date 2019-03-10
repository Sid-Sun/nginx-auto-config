import os


def getServerName():
    print("Enter the domain/subdomain name(s) (separated by space and without the ending semicolon)")
    serverName = input()
    try:
        serverName = str(serverName)
        return serverName
    except Exception:
        print("Domain MUST be a string, try again")
        serverName = getServerName()
    return serverName


def getRoot():
    print("Enter the path where the files are (root path for virtual server)")
    rootPath = input()
    try:
        rootPath = str(rootPath)
        return rootPath
    except Exception:
        print("root path MUST be a string, try again")
        rootPath = getServerName()
    return rootPath


def getURL(option):
    if option == 4:
        print("Enter the resource to proxy (EX: http://127.0.0.1:8000 or http://sidsun.com)")
    elif option == 6:
        print(
            "Enter the resource to redirect all requests to.(EX: http://sidsun.com$request_uri) (Add $request_uri if needed, it'll NOT be automatically done.")
    directPATH = input()
    try:
        rootPath = str(directPATH)
        return directPATH
    except Exception:
        print("root path MUST be a string, try again")
        directPATH = getServerName()
    return directPATH


def takeInput():
    print(
        "\n1: Create static site config (with index).\n2: Create config to host files (w/o index)\n3: Create config for Angular/Vue production site with routing\n4: Proxy pass requests to a port or a site\n5: Serve a PHP site with fastcgi and php-fpm\n6: Permanent URL redirection to someplace else.\n7: Configure to forward all HTTP requests to HTTPS\n8: Exit");
    print("\nWhat do you want to do?:")
    option = input()
    try:
        option = int(option)
        if option == 8:
            exit(0)
        else:
            os.mkdir('/etc/nginx/conf.d', 755)
            if option != 7:
                serverName = getServerName()
                if option != 4 or option != 6:
                    rootPATH = getRoot()
                else:
                    directURL = getURL(option)
                configFileName = '/etc/nginx/conf.d/' + serverName + '.conf'
                configFile = open(configFileName, 'w')
                configFile.close()
                configFile = open(configFileName, 'a')
                configFileString = 'server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n'
                if option == 1:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    location / {\n        root ' + rootPATH + ';\n        index index.html\n    }\n}\n'
                elif option == 2:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    location / {\n        root ' + rootPATH + ';\n    }\n}\n'
                elif option == 3:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    root ' + rootPATH + ';\n    index index.html;\n    location / {\n        try_files $uri $uri/ @rewrites;\n    }\n    location @rewrites {\n        rewrite ^(.+)$ /index.html last;\n    }\n}\n'
                elif option == 4:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    location / {\n        proxy_pass ' + directURL + ';\n        proxy_read_timeout  90;\n    }\n}\n'
                elif option == 5:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    root ' + rootPATH + ';\n    index index.php;\n    location / {\n        try_files $uri $uri/ =404;\n        autoindex  on;\n        autoindex_exact_size off;\n        autoindex_localtime on;\n    }\n    location ~* \.php$ {\n        include snippets/fastcgi-php.conf;\n        fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;\n    }\n}\n'
                elif option == 6:
                    configFileString = configFileString + '    ssl_certificate /etc/certbot/live/' + serverName + '/fullchain.pem;\n    ssl_certificate /etc/certbot/live/' + serverName + '/privkey.pem;\n' + '    server_name ' + serverName + ';\n    return 301 ' + directURL + ';\n}\n'
                else:
                    print('Excuse me? Run the script again and enter a PROPER number in the defined range.')
                    exit(1)
                configFile.write(configFileString)
                configFile.close()
            else:
                configFileName = '/etc/nginx/conf.d/default.conf'
                configFile = open(configFileName, 'w')
                configFile.close()
                configFile = open(configFileName, 'a')
                configFileString = 'server {\n    listen 80 default_server;\n    listen [::]:80 default_server;\n    access_log off;\n    error_log /dev/null crit;\n    server_name _;\n    return 301 https://$host$request_uri;\n}\n'
                configFile.write(configFileString)
    except Exception:
        print("Uhh, you need to type what you want to do.")
        takeInput()


def main():
    print("--------------------------------------------------------------")
    print("Script to Automate nginx virtual server creation by: Sid Sun.")
    print("Licensed under the MIT License.")
    print("By using this script, you agree to abide by the MIT License..")
    print("Copyright (c) 2019 Sidharth Soni (Sid Sun).")
    print("--------------------------------------------------------------\n")
    print("Let's  get started!")
    takeInput()


main()
