package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

const version string = "6.0.0" // Program Version

type service struct {
	selection  int
	domains    string
	root       string
	url        string
	port       int
	additional additions
}

type additions struct {
	addHSTSConfig     bool
	addSecurityConfig bool
	makeDefaultServer bool
	addCachingConfig  bool
	maxCacheAge       string
}

var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

func main() {
	var server service

	if len(os.Args) > 1 {
		pArg := os.Args[1]
		if pArg == "-h" || pArg == "-help" || pArg == "--help" {
			fmt.Println("nginx-auto-config is a program which allows you to create configurations for the nginx web server using a number of presets interactively\nLicensed under the MIT license, created by Sidharth Soni (Sid Sun)\nYou can find the source code at: https://github.com/Sid-Sun/nginx-auto-config")
		} else if pArg == "-v" || pArg == "-version" || pArg == "--version" {
			fmt.Println(version)
		} else {
			fmt.Printf(
				"Unknown option: %s\n"+
					"Run with -h, -help or --help to get help\n"+
					"-v, -version or --version to get program version\n"+
					"Or without any argumets to launch the program interactively\n", pArg)
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("An interactive program to Automate nginx virtual server creation by Sid Sun")
	fmt.Println("Licensed under the MIT License")
	fmt.Println("By using this program, you agree to abide by the MIT License")
	fmt.Println("Copyright (c) 2019 Sidharth Soni (Sid Sun)")
	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Printf("Let's  get started!\n\n")

	testWritePermissions() // Test writing permissions before proceeding further, Functions exits the program if permissions lack
	server.selection = takeInput()
	server.port = 443

	if inRange(server.selection, []int{1, 2, 3, 4, 5, 6, 7}) {
		fmt.Println("Enter the domain/sub-domain name(s) (separated by space and without ending semicolon)")
		_, _ = cyan.Print("Server Names: ")
		inputConfig := newInputConfig(false, false, "Server Names: ")
		server.domains = getInput(inputConfig)
	}

	if inRange(server.selection, []int{1, 2, 3, 4}) {
		fmt.Println("Enter the path where the files are (root path for virtual server)")
		_, _ = cyan.Print("Root path: ")
		server.root = getRootPath()
		fmt.Print("Do you want to leverage caching?")
		_, _ = cyan.Print("\nSetup Caching (Y[es]/n[o]): ")
		server.additional.addCachingConfig = getConsent(true)
		if server.additional.addCachingConfig {
			_, _ = cyan.Print("Set cache expiry [1m/4h/2d/1y] (empty for 6h): ")
			inputConfig := newInputConfig(true, true, "")
			server.additional.maxCacheAge = getInput(inputConfig)
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
		inputConfig := newInputConfig(false, true, "Root path: ")
		server.url = getInput(inputConfig)
	}

	if server.selection == 7 {
		fmt.Println("Enter the port number the virtual server should listen to")
		_, _ = cyan.Print("Port: ")
		server.port = getInt(false, "Port: ")
	}

	if server.selection == 8 {
		server.additional.makeDefaultServer = true
		server.domains = "_"
		server.port = 80
	}

	if server.selection == 9 {
		os.Exit(0)
	}

	fmt.Print("Do you want the virtual server to send HSTS preload header with the response?")
	_, _ = cyan.Print("\nSend HSTS Preload header (Y[es]/n[o]): ")
	server.additional.addHSTSConfig = getConsent(true)

	fmt.Print("Do you want to add additional security options to the config? (should not but may break the config)")
	_, _ = cyan.Print("\nAdd security config (Y[es]/n[o]): ")
	server.additional.addSecurityConfig = getConsent(true)

	fileName, fileContents := PrepareServiceFileContents(server)
	fmt.Print(fileContents)
	_, _ = cyan.Print("Is this correct? (Y[es]/n[o]): ")

	if getConsent(true) {
		writeContentToFile(fileName+".nginxAutoConfig.conf", fileContents)
		if server.port == 443 {
			printCautionSSL()
		}
	}
}

func PrepareServiceFileContents(server service) (string, string) {
	fileName := strings.Fields(server.domains)[0]
	output := "server {"
	newLine := "\n    "
	ipv4listenString := "listen " + strconv.Itoa(server.port)
	ipv6listenString := "listen [::]:" + strconv.Itoa(server.port)
	if server.additional.makeDefaultServer {
		ipv4listenString += " default_server"
		ipv6listenString += " default_server"
	}
	if server.port == 443 {
		ipv4listenString += " ssl"
		ipv6listenString += " ssl"
		ipv4listenString = "#" + ipv4listenString
		ipv6listenString = "#" + ipv6listenString
	}
	ipv4listenString += " http2;"
	ipv6listenString += " http2;"
	output += newLine + ipv4listenString
	output += newLine + ipv6listenString
	output += newLine + "server_name " + server.domains + ";"
	output += newLine + "access_log off;"
	output += newLine + "error_log /dev/null crit;"
	if server.port == 443 {
		output += newLine + "#ssl_protocols TLSv1.2 TLSv1.3;"
		output += newLine + "#ssl_certificate /etc/letsencrypt/live/" + fileName + "/fullchain.pem;"
		output += newLine + "#ssl_certificate_key /etc/letsencrypt/live/" + fileName + "/privkey.pem;"
	}
	if server.additional.addHSTSConfig {
		output += newLine + "#Send HSTS header"
		output += newLine + "add_header Strict-Transport-Security \"max-age=31536000; includeSubDomains; preload\";"
	}
	switch server.selection {
	case 1:
		output += newLine + "root " + server.root + ";"
		output += newLine + "location / {"
		output += newLine + "    index index.html;"
		output += newLine + "}"
	case 2:
		output += newLine + "location / {"
		output += newLine + "    root " + server.root + ";"
		output += newLine + "}"
	case 3:
		output += newLine + "root " + server.root + ";"
		output += newLine + "index index.html;"
		output += newLine + "location / {"
		output += newLine + "    try_files $uri $uri/ @rewrites;"
		output += newLine + "}"
		output += newLine + "location @rewrites {"
		output += newLine + "    rewrite ^(.+)$ /index.html last;"
		output += newLine + "}"
	case 4:
		output += newLine + "root " + server.root + ";"
		output += newLine + "index index.php;"
		output += newLine + "location / {"
		output += newLine + "    try_files $uri $uri/ =404;"
		output += newLine + "    autoindex  on;"
		output += newLine + "    autoindex_exact_size off;"
		output += newLine + "    autoindex_localtime on;"
		output += newLine + "}"
		output += newLine + "location ~* \\.php$ {"
		output += newLine + "    include snippets/fastcgi-php.conf;"
		output += newLine + "    fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;"
		output += newLine + "}"
	case 5, 7:
		output += newLine + "location / {"
		output += newLine + "    proxy_pass " + server.url + ";"
		output += newLine + "    proxy_read_timeout  90;"
		output += newLine + "}"
	case 6:
		output += newLine + "return 308 " + server.url + ";"
	case 8:
		fileName = "default"
		output += newLine + "return 308 https://$host$request_uri;"
	}
	if server.additional.addSecurityConfig {
		output += newLine + "#Turn off nginx version number displayed on all auto generated error pages"
		output += newLine + "server_tokens off;"
		output += newLine + "#Controlling Buffer Overflow Attacks"
		output += newLine + "#Start: Size Limits & Buffer Overflows"
		output += newLine + "client_body_buffer_size 1K;"
		output += newLine + "client_header_buffer_size 1k;"
		output += newLine + "client_max_body_size 1k;"
		output += newLine + "large_client_header_buffers 2 1k;"
		output += newLine + "#END: Size Limits & Buffer Overflows"
		output += newLine + "#Start: Timeouts"
		output += newLine + "client_body_timeout 10;"
		output += newLine + "client_header_timeout 10;"
		output += newLine + "keepalive_timeout 5 5;"
		output += newLine + "send_timeout 10;"
		output += newLine + "#End: Timeout"
		output += newLine + "#Avoid clickjacking"
		output += newLine + "add_header X-Frame-Options SAMEORIGIN;"
		output += newLine + "#Disable content-type sniffing on some browsers"
		output += newLine + "add_header X-Content-Type-Options nosniff;"
		output += newLine + "#Enable the Cross-site scripting (XSS) filter"
		output += newLine + "add_header X-XSS-Protection \"1; mode=block\";"
	}
	if server.additional.addCachingConfig {
		output += newLine + "location ~* \\.(js|css|json|png|jpg|jpeg|gif|ico)$ {"
		if server.additional.maxCacheAge == "" {
			server.additional.maxCacheAge = "6h"
		}
		output += newLine + "    expires " + server.additional.maxCacheAge + ";"
		output += newLine + "    add_header Cache-Control \"public, no-transform\";"
		output += newLine + "}"
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
	input := getInt(false, "What do you want to do: ")
	if input > 9 || input <= 0 {
		fmt.Println("Enter a valid number.")
		return takeInput()
	}
	return input
}
