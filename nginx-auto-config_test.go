package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareServiceFileContents(t *testing.T) {
	testCases := []struct {
		name                 string
		exec                 func() (string, string)
		expectedFileName     string
		expectedFileContents string
	}{
		{
			name: "test create service for static website hosting with HSTS, Security",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     true,
					addSecurityConfig: true,
					makeDefaultServer: false,
				}
				service := service{
					selection:  1,
					domains:    "sidsun.com cdn.sidsun.com",
					root:       "/srv/www/sid",
					url:        "",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "sidsun.com",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name sidsun.com cdn.sidsun.com;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/sidsun.com/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/sidsun.com/privkey.pem;
    #Send HSTS header
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";
    root /srv/www/sid;
    location / {
        index index.html;
    }
    #Turn off nginx version number displayed on all auto generated error pages
    server_tokens off;
    #Controlling Buffer Overflow Attacks
    #Start: Size Limits & Buffer Overflows
    client_body_buffer_size 1K;
    client_header_buffer_size 1k;
    client_max_body_size 1k;
    large_client_header_buffers 2 1k;
    #END: Size Limits & Buffer Overflows
    #Start: Timeouts
    client_body_timeout 10;
    client_header_timeout 10;
    keepalive_timeout 5 5;
    send_timeout 10;
    #End: Timeout
    #Avoid clickjacking
    add_header X-Frame-Options SAMEORIGIN;
    #Disable content-type sniffing on some browsers
    add_header X-Content-Type-Options nosniff;
    #Enable the Cross-site scripting (XSS) filter
    add_header X-XSS-Protection "1; mode=block";
}
`,
		},
		{
			name: "test create service for hosting files without index with Security",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     false,
					addSecurityConfig: true,
					makeDefaultServer: false,
				}
				service := service{
					selection:  2,
					domains:    "sulabs.org",
					root:       "/srv/www/su",
					url:        "",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "sulabs.org",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name sulabs.org;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/sulabs.org/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/sulabs.org/privkey.pem;
    location / {
        root /srv/www/su;
    }
    #Turn off nginx version number displayed on all auto generated error pages
    server_tokens off;
    #Controlling Buffer Overflow Attacks
    #Start: Size Limits & Buffer Overflows
    client_body_buffer_size 1K;
    client_header_buffer_size 1k;
    client_max_body_size 1k;
    large_client_header_buffers 2 1k;
    #END: Size Limits & Buffer Overflows
    #Start: Timeouts
    client_body_timeout 10;
    client_header_timeout 10;
    keepalive_timeout 5 5;
    send_timeout 10;
    #End: Timeout
    #Avoid clickjacking
    add_header X-Frame-Options SAMEORIGIN;
    #Disable content-type sniffing on some browsers
    add_header X-Content-Type-Options nosniff;
    #Enable the Cross-site scripting (XSS) filter
    add_header X-XSS-Protection "1; mode=block";
}
`,
		},
		{
			name: "test create service for routed webapp hosting with HSTS",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     true,
					addSecurityConfig: false,
					makeDefaultServer: false,
				}
				service := service{
					selection:  3,
					domains:    "encrypt.ml",
					root:       "/srv/www/encrypt",
					url:        "",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "encrypt.ml",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name encrypt.ml;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/encrypt.ml/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/encrypt.ml/privkey.pem;
    #Send HSTS header
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";
    root /srv/www/encrypt;
    index index.html;
    location / {
        try_files $uri $uri/ @rewrites;
    }
    location @rewrites {
        rewrite ^(.+)$ /index.html last;
    }
}
`,
		},
		{
			name: "test create service for basic PHP Website hosting",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     false,
					addSecurityConfig: false,
					makeDefaultServer: false,
				}
				service := service{
					selection:  4,
					domains:    "strangebits.co.in beta.strangebits.co.in",
					root:       "/srv/www/strange",
					url:        "",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "strangebits.co.in",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name strangebits.co.in beta.strangebits.co.in;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/strangebits.co.in/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/strangebits.co.in/privkey.pem;
    root /srv/www/strange;
    index index.php;
    location / {
        try_files $uri $uri/ =404;
        autoindex  on;
        autoindex_exact_size off;
        autoindex_localtime on;
    }
    location ~* \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass  unix:/var/run/php/php7.2-fpm.sock;
    }
}
`,
		},
		{
			name: "test create service for Proxy requests with HSTS and security",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     true,
					addSecurityConfig: true,
					makeDefaultServer: false,
				}
				service := service{
					selection:  5,
					domains:    "sulabs.ml writewith.me",
					root:       "",
					url:        "https://blog.sidsun.com",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "sulabs.ml",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name sulabs.ml writewith.me;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/sulabs.ml/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/sulabs.ml/privkey.pem;
    #Send HSTS header
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";
    location / {
        proxy_pass https://blog.sidsun.com;
        proxy_read_timeout  90;
    }
    #Turn off nginx version number displayed on all auto generated error pages
    server_tokens off;
    #Controlling Buffer Overflow Attacks
    #Start: Size Limits & Buffer Overflows
    client_body_buffer_size 1K;
    client_header_buffer_size 1k;
    client_max_body_size 1k;
    large_client_header_buffers 2 1k;
    #END: Size Limits & Buffer Overflows
    #Start: Timeouts
    client_body_timeout 10;
    client_header_timeout 10;
    keepalive_timeout 5 5;
    send_timeout 10;
    #End: Timeout
    #Avoid clickjacking
    add_header X-Frame-Options SAMEORIGIN;
    #Disable content-type sniffing on some browsers
    add_header X-Content-Type-Options nosniff;
    #Enable the Cross-site scripting (XSS) filter
    add_header X-XSS-Protection "1; mode=block";
}
`,
		},
		{
			name: "test create service for Permanent URL redirection with security",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     false,
					addSecurityConfig: true,
					makeDefaultServer: false,
				}
				service := service{
					selection:  6,
					domains:    "strangebits.co.in readwith.me",
					root:       "",
					url:        "http://blog.sidsun.com$request_uri",
					port:       443,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "strangebits.co.in",
			expectedFileContents: `server {
    listen 443;
    listen [::]:443;
    server_name strangebits.co.in readwith.me;
    access_log off;
    error_log /dev/null crit;
    #ssl on;
    #ssl_certificate /etc/letsencrypt/live/strangebits.co.in/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/strangebits.co.in/privkey.pem;
    return 308 http://blog.sidsun.com$request_uri;
    #Turn off nginx version number displayed on all auto generated error pages
    server_tokens off;
    #Controlling Buffer Overflow Attacks
    #Start: Size Limits & Buffer Overflows
    client_body_buffer_size 1K;
    client_header_buffer_size 1k;
    client_max_body_size 1k;
    large_client_header_buffers 2 1k;
    #END: Size Limits & Buffer Overflows
    #Start: Timeouts
    client_body_timeout 10;
    client_header_timeout 10;
    keepalive_timeout 5 5;
    send_timeout 10;
    #End: Timeout
    #Avoid clickjacking
    add_header X-Frame-Options SAMEORIGIN;
    #Disable content-type sniffing on some browsers
    add_header X-Content-Type-Options nosniff;
    #Enable the Cross-site scripting (XSS) filter
    add_header X-XSS-Protection "1; mode=block";
}
`,
		},
		{
			name: "test create service for basic proxy with custom port",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     false,
					addSecurityConfig: false,
					makeDefaultServer: false,
				}
				service := service{
					selection:  7,
					domains:    "_",
					root:       "",
					url:        "http://127.0.0.1:5000",
					port:       4321,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "_",
			expectedFileContents: `server {
    listen 4321;
    listen [::]:4321;
    server_name _;
    access_log off;
    error_log /dev/null crit;
    location / {
        proxy_pass http://127.0.0.1:5000;
        proxy_read_timeout  90;
    }
}
`,
		},
		{
			name: "test create default service without HSTS	and security config",
			exec: func() (string, string) {
				additions := additions{
					addHSTSConfig:     false,
					addSecurityConfig: false,
					makeDefaultServer: true,
				}
				service := service{
					selection:  8,
					domains:    "_",
					root:       "",
					url:        "",
					port:       80,
					additional: additions,
				}
				return PrepareServiceFileContents(service)
			},
			expectedFileName: "default",
			expectedFileContents: `server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name _;
    access_log off;
    error_log /dev/null crit;
    return 308 https://$host$request_uri;
}
`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fileName, fileContents := testCase.exec()
			assert.Equal(t, testCase.expectedFileName, fileName)
			assert.Equal(t, testCase.expectedFileContents, fileContents)
		})
	}
}
