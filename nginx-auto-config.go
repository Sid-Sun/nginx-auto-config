package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

type service struct {
	selection  int
	domains    string
	root       string
	verifyRoot bool
	url        string
	port       int
	additional additions
}

type additions struct {
	addHSTSConfig     bool
	addSecurityConfig bool
	makeDefaultServer bool
}

var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

func main() {
	var server service
	server.verifyRoot = true

	if len(os.Args) > 1 {
		pArgs := os.Args[1]
		if pArgs == "-h" || pArgs == "-help" || pArgs == "--help" {
			fmt.Println("nginx-auto-config is a program which allows you to create configurations for the nginx web server using a number of presets interactively\nLicensed under the MIT license, created by Sidharth Soni (Sid Sun)\nYou can find the source code at: https://github.com/Sid-Sun/nginx-auto-config")
			os.Exit(0)
		} else if pArgs == "version" {
			fmt.Printf("3.2\n")
			os.Exit(0)
		} else if pArgs == "-skiproot" || pArgs == "-s" || pArgs == "--skiproot" {
			server.verifyRoot = false
		} else {
			fmt.Printf("Unknown option(s) %s, run with -h, -help or --help to get help, -s, -skiproot or --skiproot to skip server root directory validation or without any argumets to launch the program\n", pArgs)
			os.Exit(1)
		}
	}
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("An interactive program to Automate nginx virtual server creation by Sid Sun")
	fmt.Println("Licensed under the MIT License")
	fmt.Println("By using this program, you agree to abide by the MIT License")
	fmt.Println("Copyright (c) 2019 Sidharth Soni (Sid Sun)")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Printf("Let's  get started!\n\n")
	testWritePermissions()
	server.selection = takeInput()
	server.port = 443
	if inRange(server.selection, []int{1, 2, 3, 4, 5, 6, 7}) {
		fmt.Println("Enter the domain/sub-domain name(s) (separated by space and without ending semicolon)")
		_, _ = cyan.Print("Server Names: ")
		server.domains = getInput(false, false)
	}
	if inRange(server.selection, []int{1, 2, 3, 4}) {
		fmt.Println("Enter the path where the files are (root path for virtual server)")
		_, _ = cyan.Print("Root path: ")
		rootPath := getInput(false, false)
		if server.verifyRoot && !dirExists(rootPath) {
			_, _ = red.Printf("Server root directory: '%v' is non existent. Please try again\n", rootPath)
			os.Exit(1)
		} else {
			server.root = rootPath
		}
	}
	if inRange(server.selection, []int{5, 6, 7}) {
		if server.selection == 6 {
			fmt.Println("Enter the resource to redirect all requests to.(EX: http://sidsun.com$request_uri) (Add $request_uri if needed, it'll NOT be automatically done)")
			_, _ = cyan.Print("Redirect URL: ")
		} else {
			fmt.Println("Enter the resource to proxy (EX: http://127.0.0.1:8000 or http://sidsun.com)")
			_, _ = cyan.Print("Resource to proxy: ")
		}
		server.url = getInput(false, true)
	}
	if server.selection == 7 {
		fmt.Println("Enter the port number the virtual server should listen to")
		_, _ = cyan.Print("Port: ")
		server.port = getInt()
	}
	if server.selection == 8 {
		server.additional.makeDefaultServer = true
		server.domains = "_"
		server.port = 80
	}
	if server.selection == 9 {
		os.Exit(1)
	}
	fmt.Print("Do you want the virtual server to send HSTS preload header with the response? (Y[es]/N[o]): ")
	_, _ = cyan.Print("\nSend HSTS Preload header (Y[es]/N[o]): ")
	server.additional.addHSTSConfig = getConsent()
	fmt.Print("Do you want to add additional security options to the config? (should not but may break the config) (Y[es]/N[o]): ")
	_, _ = cyan.Print("\nAdd security config (Y[es]/N[o]): ")
	server.additional.addSecurityConfig = getConsent()
	fileName, fileContents := PrepareServiceFileContents(server)
	fmt.Print(fileContents)
	_, _ = cyan.Print("Is this correct? (Y[es]/N[o]): ")
	if getConsent() {
		writeContentToFile(fileName+".nginxAutoConfig.conf", fileContents)
		if server.port == 443 {
			printCautionSSL()
		}
		os.Exit(0)
	}

	os.Exit(0)
}

func PrepareServiceFileContents(server service) (string, string) {
	fileName := strings.Fields(server.domains)[0]
	output := "server {"
	if server.additional.makeDefaultServer {
		output += "\n    listen " + strconv.Itoa(server.port) + " default_server;\n    listen [::]:" + strconv.Itoa(server.port) + " default_server;"
	} else {
		output += "\n    listen " + strconv.Itoa(server.port) + ";\n    listen [::]:" + strconv.Itoa(server.port) + ";"
	}
	output += "\n    server_name " + server.domains + ";"
	output += "\n    access_log off;"
	output += "\n    error_log /dev/null crit;"
	if server.port == 443 {
		output += "\n    #ssl on;"
		output += "\n    #ssl_certificate /etc/letsencrypt/live/" + fileName + "/fullchain.pem;"
		output += "\n    #ssl_certificate_key /etc/letsencrypt/live/" + fileName + "/privkey.pem;"
	}
	if server.additional.addHSTSConfig {
		output += "\n    #Send HSTS header"
		output += "\n    add_header Strict-Transport-Security \"max-age=31536000; includeSubDomains; preload\";"
	}
	switch server.selection {
	case 1:
		output += "\n    location / {"
		output += "\n        index index.html;"
		output += "\n        root " + server.root + ";"
		output += "\n    }"
	case 2:
		output += "\n    location / {"
		output += "\n        root " + server.root + ";"
		output += "\n    }"
	case 3:
		output += "\n    root " + server.root + ";"
		output += "\n    index index.html;"
		output += "\n    location / {"
		output += "\n        try_files $uri $uri/ @rewrites;"
		output += "\n    }"
		output += "\n    location @rewrites {"
		output += "\n        rewrite ^(.+)$ /index.html last;"
		output += "\n    }"
	case 4:
		output += "\n    root " + server.root + ";"
		output += "\n    index index.php;"
		output += "\n    location / {"
		output += "\n        try_files $uri $uri/ =404;"
		output += "\n        autoindex  on;"
		output += "\n        autoindex_exact_size off;"
		output += "\n        autoindex_localtime on;"
		output += "\n    }"
		output += "\n    location ~* \\.php$ {"
		output += "\n        include snippets/fastcgi-php.conf;"
		output += "\n        fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;"
		output += "\n    }"
	case 5, 7:
		output += "\n    location / {"
		output += "\n        proxy_pass " + server.url + ";"
		output += "\n        proxy_read_timeout  90;"
		output += "\n    }"
	case 6:
		output += "\n    return 308 " + server.url + ";"
	case 8:
		fileName = "default"
		output += "\n    return 308 https://$host$request_uri;"
	}
	if server.additional.addSecurityConfig {
		output += "\n    #Turn off nginx version number displayed on all auto generated error pages"
		output += "\n    server_tokens off;"
		output += "\n    #Controlling Buffer Overflow Attacks"
		output += "\n    #Start: Size Limits & Buffer Overflows"
		output += "\n    client_body_buffer_size 1K;"
		output += "\n    client_header_buffer_size 1k;"
		output += "\n    client_max_body_size 1k;"
		output += "\n    large_client_header_buffers 2 1k;"
		output += "\n    #END: Size Limits & Buffer Overflows"
		output += "\n    #Start: Timeouts"
		output += "\n    client_body_timeout 10;"
		output += "\n    client_header_timeout 10;"
		output += "\n    keepalive_timeout 5 5;"
		output += "\n    send_timeout 10;"
		output += "\n    #End: Timeout"
		output += "\n    #Avoid clickjacking"
		output += "\n    add_header X-Frame-Options SAMEORIGIN;"
		output += "\n    #Disable content-type sniffing on some browsers"
		output += "\n    add_header X-Content-Type-Options nosniff;"
		output += "\n    #Enable the Cross-site scripting (XSS) filter"
		output += "\n    add_header X-XSS-Protection \"1; mode=block\";"
	}
	output += "\n}\n" //End server block and add newline at EOF
	return fileName, output
}

func takeInput() int {
	_, _ = yellow.Print("Options: \n")
	fmt.Println("(1) Static website hosting - Host a static website")
	fmt.Println("(2) Host files without index - Host files without an index")
	fmt.Println("(3) Host a routed webapp - Host a React/Angular/Vue webapp")
	fmt.Println("(4) PHP Website hosting - Host a PHP site with fastcgi and php-fpm")
	fmt.Println("(5) Proxy requests - Proxy incoming requests to a port or a website")
	fmt.Println("(6) Permanent URL redirection - Redirect all incoming requests to an address")
	fmt.Println("(7) Proxy with custom port - Proxy incoming requests at a port to an address")
	fmt.Println("(8) HTTP requests to HTTPS redirect - Redirects all incoming HTTP traffic to HTTPS (use as default config)")
	fmt.Println("(9) Exit")
	_, _ = cyan.Print("What do you want to do: ")
	input := getInt()
	if input > 9 || input <= 0 {
		fmt.Println("Enter a valid number.")
		return takeInput()
	}
	return input
}
