package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
)

const version string = "6.1.0" // Program Version

type Service struct {
	Selection  int
	Domains    string
	Root       string
	URL        string
	Port       int
	Additional Additions
}

type Additions struct {
	AddHSTSConfig     bool
	AddSecurityConfig bool
	MakeDefaultServer bool
	AddCachingConfig  bool
	MaxCacheAge       string
}

var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

func main() {
	if len(os.Args) > 1 {
		pArg := os.Args[1]
		if pArg == "-h" || pArg == "-help" || pArg == "--help" {
			fmt.Println("nginx-auto-config is a program which allows you to create configurations for the nginx web server using a number of presets interactively\nLicensed under the MIT license, created by Sidharth Soni (Sid Sun)\nYou can find the source code at: https://github.com/Sid-Sun/nginx-auto-config")
		} else if pArg == "-v" || pArg == "-version" || pArg == "--version" {
			fmt.Println(version)
		} else if fileExists(pArg) {
			serverFileConf := readFromFile(pArg)
			serviceFile := Service{}

			if err := toml.Unmarshal(serverFileConf, &serviceFile); err != nil {
				red.Println("Error occoured while reading config from", pArg, "Details: \n", err.Error())
				os.Exit(1)
			}

			fileName, fileContents := prepareServiceFileContents(serviceFile)

			fmt.Print(fileContents)
			_, _ = cyan.Print("Is this correct? (Y[es]/n[o]): ")

			if getConsent(true) {
				if err := writeContentToFile(fileName+".conf", []byte(fileContents)); err != nil {
					red.Println("Error occoured while writing config", fileName+".conf", "Details:\n", err.Error())
					os.Exit(1)
				}

				fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", fileName+".conf")
				if serviceFile.Port == 443 {
					printCautionSSL()
				}
			}

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

	serviceConfig := getDetails()

	fileName, fileContents := prepareServiceFileContents(serviceConfig)
	fmt.Print(fileContents)
	_, _ = cyan.Print("Is this correct? (Y[es]/n[o]): ")

	if getConsent(true) {
		if data, err := toml.Marshal(serviceConfig); err != nil {
			red.Println("Error occoured while converting config to toml. Details: \n", err.Error())
			os.Exit(1)
		} else {
			err = writeContentToFile(fileName+".toml", data)
			if err != nil {
				red.Println("Error occoured while writing config", fileName+".toml", "Details:\n", err.Error())
				os.Exit(1)
			}
		}

		if err := writeContentToFile(fileName+".conf", []byte(fileContents)); err != nil {
			red.Println("Error occoured while writing config", fileName+".conf", "Details:\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Wrote service details to %s, run program with %s as argument to re-generate config!\n", fileName+".toml", fileName+".toml")
		fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", fileName+".conf")

		if serviceConfig.Port == 443 {
			printCautionSSL()
		}
	}
}

func prepareServiceFileContents(server Service) (string, string) {
	fileName := strings.Fields(server.Domains)[0]
	output := "server {"
	newLine := "\n    "
	ipv4listenString := "listen " + strconv.Itoa(server.Port)
	ipv6listenString := "listen [::]:" + strconv.Itoa(server.Port)
	if server.Additional.MakeDefaultServer {
		ipv4listenString += " default_server"
		ipv6listenString += " default_server"
	}
	if server.Port == 443 {
		ipv4listenString += " ssl"
		ipv6listenString += " ssl"
		ipv4listenString = "#" + ipv4listenString
		ipv6listenString = "#" + ipv6listenString
	}
	ipv4listenString += " http2;"
	ipv6listenString += " http2;"
	output += newLine + ipv4listenString
	output += newLine + ipv6listenString
	output += newLine + "server_name " + server.Domains + ";"
	output += newLine + "access_log off;"
	output += newLine + "error_log /dev/null crit;"
	if server.Port == 443 {
		output += newLine + "#ssl_protocols TLSv1.2 TLSv1.3;"
		output += newLine + "#ssl_certificate /etc/letsencrypt/live/" + fileName + "/fullchain.pem;"
		output += newLine + "#ssl_certificate_key /etc/letsencrypt/live/" + fileName + "/privkey.pem;"
	}
	if server.Additional.AddHSTSConfig {
		output += newLine + "#Send HSTS header"
		output += newLine + "add_header Strict-Transport-Security \"max-age=31536000; includeSubDomains; preload\";"
	}
	switch server.Selection {
	case 1:
		output += newLine + "root " + server.Root + ";"
		output += newLine + "location / {"
		output += newLine + "    index index.html;"
		output += newLine + "}"
	case 2:
		output += newLine + "location / {"
		output += newLine + "    root " + server.Root + ";"
		output += newLine + "}"
	case 3:
		output += newLine + "root " + server.Root + ";"
		output += newLine + "index index.html;"
		output += newLine + "location / {"
		output += newLine + "    try_files $uri $uri/ @rewrites;"
		output += newLine + "}"
		output += newLine + "location @rewrites {"
		output += newLine + "    rewrite ^(.+)$ /index.html last;"
		output += newLine + "}"
	case 4:
		output += newLine + "root " + server.Root + ";"
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
		output += newLine + "    proxy_pass " + server.URL + ";"
		output += newLine + "    proxy_read_timeout  90;"
		output += newLine + "}"
	case 6:
		output += newLine + "return 308 " + server.URL + ";"
	case 8:
		fileName = "default"
		output += newLine + "return 308 https://$host$request_uri;"
	}
	if server.Additional.AddSecurityConfig {
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
	if server.Additional.AddCachingConfig {
		output += newLine + "location ~* \\.(js|css|json|png|jpg|jpeg|gif|ico)$ {"
		if server.Additional.MaxCacheAge == "" {
			server.Additional.MaxCacheAge = "6h"
		}
		output += newLine + "    expires " + server.Additional.MaxCacheAge + ";"
		output += newLine + "    add_header Cache-Control \"public, no-transform\";"
		output += newLine + "}"
	}
	output += "\n}\n" //End server block and add newline at EOF
	return fileName, output
}

func getDetails() Service {
	var server Service

	server.Selection = takeInput()
	server.Port = 443

	if inRange(server.Selection, []int{1, 2, 3, 4, 5, 6, 7}) {
		fmt.Println("Enter the domain/sub-domain name(s) (separated by space and without ending semicolon)")
		_, _ = cyan.Print("Server Names: ")
		inputConfig := newInputConfig(false, false, "Server Names: ")
		server.Domains = getInput(inputConfig)
	}

	if inRange(server.Selection, []int{1, 2, 3, 4}) {
		fmt.Println("Enter the path where the files are (root path for virtual server)")
		_, _ = cyan.Print("Root path: ")
		server.Root = getRootPath()
		fmt.Print("Do you want to leverage caching?")
		_, _ = cyan.Print("\nSetup Caching (Y[es]/n[o]): ")
		server.Additional.AddCachingConfig = getConsent(true)
		if server.Additional.AddCachingConfig {
			_, _ = cyan.Print("Set cache expiry [1m/4h/2d/1y] (empty for 6h): ")
			inputConfig := newInputConfig(true, true, "")
			server.Additional.MaxCacheAge = getInput(inputConfig)
		}
	}

	if inRange(server.Selection, []int{5, 6, 7}) {
		if server.Selection == 6 {
			fmt.Println("Enter the resource to redirect all requests to.(EX: http://sidsun.com$request_uri) (Add $request_uri if needed, it'll NOT be automatically done)")
			_, _ = cyan.Print("Redirect URL: ")
		} else {
			fmt.Println("Enter the resource to proxy (EX: http://127.0.0.1:8000 or http://sidsun.com)")
			_, _ = cyan.Print("Resource to proxy: ")
		}
		inputConfig := newInputConfig(false, true, "Root path: ")
		server.URL = getInput(inputConfig)
	}

	if server.Selection == 7 {
		fmt.Println("Enter the port number the virtual server should listen to")
		_, _ = cyan.Print("Port: ")
		server.Port = getInt(false, "Port: ")
	}

	if server.Selection == 8 {
		server.Additional.MakeDefaultServer = true
		server.Domains = "_"
		server.Port = 80
	}

	if server.Selection == 9 {
		os.Exit(0)
	}

	fmt.Print("Do you want the virtual server to send HSTS preload header with the response?")
	_, _ = cyan.Print("\nSend HSTS Preload header (Y[es]/n[o]): ")
	server.Additional.AddHSTSConfig = getConsent(true)

	fmt.Print("Do you want to add additional security options to the config? (should not but may break the config)")
	_, _ = cyan.Print("\nAdd security config (Y[es]/n[o]): ")
	server.Additional.AddSecurityConfig = getConsent(true)

	return server
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
