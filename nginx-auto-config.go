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
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help" {
			fmt.Println("nginx-auto-config is a program which allows you to create configurations for the nginx web server using a number of presets interactively\nLicensed under the MIT license, created by Sidharth Soni (Sid Sun)\nYou can find the source code at: https://github.com/Sid-Sun/nginx-auto-config")
		} else {
			fmt.Printf("Unknown option(s) %s, run with -h, -help or --help to get help or without any argumets to launch the program\n",os.Args[1])
		}
		return
	}
	var input uint
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("An interactive program to Automate nginx virtual server creation by Sid Sun.")
	fmt.Println("Licensed under the MIT License.")
	fmt.Println("By using this program, you agree to abide by the MIT License.")
	fmt.Println("Copyright (c) 2019 Sidharth Soni (Sid Sun).")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Printf("Let's  get started!\n\n")
	testWritePermissions()
	input = takeInput()
	var configFileName string
	var configFileContents string
	switch input {
	case 1:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        root " + rootPATH + ";\n        index index.html;\n    }\n"
	case 2:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        root " + rootPATH + ";\n    }\n"
	case 3:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    root " + rootPATH + ";\n    index index.html;\n    location / {\n        try_files $uri $uri/ @rewrites;\n    }\n    location @rewrites {\n        rewrite ^(.+)$ /index.html last;\n    }\n"
	case 4:
		serverName := getServerName()
		directURL := getURL(input)
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    location / {\n        proxy_pass " + directURL + ";\n        proxy_read_timeout  90;\n    }\n"
	case 5:
		serverName := getServerName()
		rootPATH := getRoot()
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    root " + rootPATH + ";\n    index index.php;\n    location / {\n        try_files $uri $uri/ =404;\n        autoindex  on;\n        autoindex_exact_size off;\n        autoindex_localtime on;\n    }\n    location ~* \\.php$ {\n        include snippets/fastcgi-php.conf;\n        fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;\n    }\n"
	case 6:
		serverName := getServerName()
		directURL := getURL(input)
		configFileName = createConfigFile(serverName)
		configFileContents := "server {\n    listen 443;\n    listen [::]:443;\n    ssl on;\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    ssl_certificate /etc/letsencrypt/live/" + serverName + "/fullchain.pem;\n    ssl_certificate_key /etc/letsencrypt/live/" + serverName + "/privkey.pem;\n" + "    server_name " + serverName + ";\n    return 301 " + directURL + ";\n"
	case 7:
		serverName := "default"
		configFileName = createConfigFile(serverName)
		configFileContents = "server {\n    listen 80 default_server;\n    listen [::]:80 default_server;\n    access_log off;\n    error_log /dev/null crit;\n    server_name _;\n    return 301 https://$host$request_uri;\n"
	case 8:
		directURL := getURL(input)
		portNumber := strconv.Itoa(getListeningPort())
		serverName := getVirtualServerAlias()
		configFileName = createConfigFile(serverName)
		configFileContents = "server {\n    listen " + portNumber + ";\n    listen [::]:" + portNumber + ";\n    access_log off;\n    error_log /dev/null crit;\n"
		configFileContents = configFileContents + "    server_name _;\n    location / {\n        proxy_pass " + directURL + ";\n        proxy_read_timeout  90;\n    }\n"
	case 9:
		return
	default:
		os.Exit(1)
	}
	fmt.Print("Do you want to add additional security options to the config? (should not but may break the config)?\n(yes/no): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	addSecurityOptions := scanner.Text()
	if addSecurityOptions == "yes" {
		configFileContents = configFileContents + "    #Turn off nginx version number displayed on all auto generated error pages\n    server_tokens off;\n    #Controlling Buffer Overflow Attacks\n    #Start: Size Limits & Buffer Overflows\n    client_body_buffer_size 1K;\n    client_header_buffer_size 1k;\n    client_max_body_size 1k;\n    large_client_header_buffers 2 1k;\n    #END: Size Limits & Buffer Overflows\n    #Start: Timeouts\n    client_body_timeout 10;\n    client_header_timeout 10;\n    keepalive_timeout 5 5;\n    send_timeout 10;\n    #End: Timeout\n    #Avoid clickjacking\n    add_header X-Frame-Options SAMEORIGIN;\n    #Disable content-type sniffing on some browsers\n    add_header X-Content-Type-Options nosniff;\n    #Enable the Cross-site scripting (XSS) filter\n    add_header X-XSS-Protection \"1; mode=block\";\n"
	}
	configFileContents = configFileContents + "}\n"
	writeContentToFile(configFileName, configFileContents)
}

func takeInput() uint {
	var input int
	fmt.Println("What do you want to do?")
	fmt.Printf("\n1: Create static site config (with index).\n2: Create config to host files (w/o index)\n3: Create config for Angular/Vue production site with routing\n4: Proxy pass requests to a port or a site\n5: Serve a PHP site with fastcgi and php-fpm\n6: Permanent URL redirection to someplace else.\n7: Configure to forward all HTTP requests to HTTPS\n8: Port forward without hostname with custom port numbers\n9: Exit\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scannedText := scanner.Text()
	input, err := strconv.Atoi(scannedText)
	if err != nil {
		fmt.Println("Something went wrong, please try again.")
		return takeInput()
	}
	if uint(input) > 9 {
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

func getListeningPort() int {
	fmt.Println("Enter the port number the virtual server should listen to")
	var input int
	_, _ = fmt.Scanf("%d", &input)
	if input == 0 {
		return getListeningPort()
	}
	return input
}

func getVirtualServerAlias() string {
	fmt.Println("Enter an alias for the virtual server (used for file name referencing)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getURL(option uint) string {
	if option == 4 || option == 8 {
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

func writeContentToFile(fileName string, fileContents string) {
	testWritePermissions()
	err := ioutil.WriteFile(fileName, []byte(fileContents), 0644)
	if err != nil {
		fmt.Println("Something went wrong, please send the log below to Sid Sun.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", fileName)
}
