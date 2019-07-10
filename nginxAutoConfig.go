package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	var input uint
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("An interactive program to Automate nginx virtual server creation by Sid Sun.")
	fmt.Println("Licensed under the MIT License.")
	fmt.Println("By using this script, you agree to abide by the MIT License.")
	fmt.Println("Copyright (c) 2019 Sidharth Soni (Sid Sun).")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Printf("Let's  get started!\n\n")
	testWritePermissions()
	input = takeInput()
	switch input {
	case 1:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        root " + rootPATH + ";\n        index index.html;\n    }\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 2:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        root " + rootPATH + ";\n    }\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 3:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    root " + rootPATH + ";\n    index index.html;\n    location / {\n        try_files $uri $uri/ @rewrites;\n    }\n    location @rewrites {\n        rewrite ^(.+)$ /index.html last;\n    }\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 4:
		serverName := getServerName()
		directURL := getURL(input)
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        proxy_pass " + directURL + ";\n        proxy_read_timeout  90;\n    }\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 5:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    root " + rootPATH + ";\n    index index.php;\n    location / {\n        try_files $uri $uri/ =404;\n        autoindex  on;\n        autoindex_exact_size off;\n        autoindex_localtime on;\n    }\n    location ~* \\.php$ {\n        include snippets/fastcgi-php.conf;\n        fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;\n    }\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 6:
		serverName := getServerName()
		directURL := getURL(input)
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/certbot/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/certbot/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    return 301 " + directURL + ";\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 7:
		serverName := "default"
		configFileName := createConfigFile(serverName)
		configFileContents := "server {\n    listen 80 default_server;\n    listen [::]:80 default_server;\n    access_log off;\n    error_log /dev/null crit;\n    server_name _;\n    return 301 https://$host$request_uri;\n}\n"
		if writeContentToFile(configFileName, configFileContents) {
			fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", configFileName)
		}
	case 8:
		os.Exit(0)
	default:
		os.Exit(1)
	}
}

func takeInput() uint {
	var input int
	fmt.Println("What do you want to do?")
	fmt.Printf("\n1: Create static site config (with index).\n2: Create config to host files (w/o index)\n3: Create config for Angular/Vue production site with routing\n4: Proxy pass requests to a port or a site\n5: Serve a PHP site with fastcgi and php-fpm\n6: Permanent URL redirection to someplace else.\n7: Configure to forward all HTTP requests to HTTPS\n8: Exit\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scannedText := scanner.Text()
	input, err := strconv.Atoi(scannedText)
	if err != nil {
		fmt.Println("Something went wrong, please try again.")
		return takeInput()
	}
	if uint(input) > 8 {
		fmt.Println("Enter a valid number.")
		return takeInput()
	}
	return uint(input)
}

func getServerName() string {
	fmt.Println("Enter the domain/subdomain name(s) (separated by space and without the ending semicolon")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getRoot() string {
	fmt.Println("Enter the path where the files are (root path for virtual server)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getURL(option uint) string {
	if option == 4 {
		fmt.Println("Enter the resource to proxy (EX: http://127.0.0.1:8000 or http://sidsun.com)")
	} else if option == 6 {
		fmt.Println("Enter the resource to redirect all requests to.(EX: http://sidsun.com$request_uri) (Add $request_uri if needed, it'll NOT be automatically done.")
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func testWritePermissions() {
	newFile, err := os.Create("nginxAutoConfig.test.txt")
	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: Write permission denied, please cd into a workable dir, exiting.")
			os.Exit(1)
		}
	} else {
		_ = os.Remove("nginxAutoConfig.test.txt")
		_ = newFile.Close()
	}
}

func createConfigFile(serverName string) string {
	configFileName := serverName + ".nginxAutoConfig.conf"
	configFile, _ := os.Create(configFileName)
	_ = configFile.Close()
	return string(serverName + ".nginxAutoConfig.conf")
}

func writeContentToFile(fileName string, fileContents string) bool {
	err := ioutil.WriteFile(fileName, []byte(fileContents), 0644)
	if err != nil {
		fmt.Println("Something went wrong, please send the log below to Sid Sun.")
		fmt.Println(err)
		return false
	}
	return true
}
