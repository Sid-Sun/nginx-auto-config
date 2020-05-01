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
				additions := Additions{
					AddHSTSConfig:     true,
					AddSecurityConfig: true,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  1,
					Domains:    "sidsun.com cdn.sidsun.com",
					Root:       "/srv/www/sid",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "sidsun.com",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name sidsun.com cdn.sidsun.com;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: true,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  2,
					Domains:    "sulabs.org",
					Root:       "/srv/www/su",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "sulabs.org",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name sulabs.org;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     true,
					AddSecurityConfig: false,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  3,
					Domains:    "encrypt.ml",
					Root:       "/srv/www/encrypt",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "encrypt.ml",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name encrypt.ml;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: false,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  4,
					Domains:    "strangebits.co.in beta.strangebits.co.in",
					Root:       "/srv/www/strange",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "strangebits.co.in",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name strangebits.co.in beta.strangebits.co.in;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     true,
					AddSecurityConfig: true,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  5,
					Domains:    "sulabs.ml writewith.me",
					Root:       "",
					URL:        "https://blog.sidsun.com",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "sulabs.ml",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name sulabs.ml writewith.me;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: true,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  6,
					Domains:    "strangebits.co.in readwith.me",
					Root:       "",
					URL:        "http://blog.sidsun.com$request_uri",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "strangebits.co.in",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name strangebits.co.in readwith.me;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
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
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: false,
					MakeDefaultServer: false,
				}
				service := Service{
					Selection:  7,
					Domains:    "_",
					Root:       "",
					URL:        "http://127.0.0.1:5000",
					Port:       4321,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "_",
			expectedFileContents: `server {
    listen 4321 http2;
    listen [::]:4321 http2;
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
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: false,
					MakeDefaultServer: true,
				}
				service := Service{
					Selection:  8,
					Domains:    "_",
					Root:       "",
					URL:        "",
					Port:       80,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "default",
			expectedFileContents: `server {
    listen 80 default_server http2;
    listen [::]:80 default_server http2;
    server_name _;
    access_log off;
    error_log /dev/null crit;
    return 308 https://$host$request_uri;
}
`,
		},
		{
			name: "test create service for static website hosting with caching",
			exec: func() (string, string) {
				additions := Additions{
					AddHSTSConfig:     false,
					AddSecurityConfig: false,
					MakeDefaultServer: false,
					AddCachingConfig:  true,
					MaxCacheAge:       "",
				}
				service := Service{
					Selection:  1,
					Domains:    "sulabs.ml encrypt.ml",
					Root:       "/srv/www/sulabs",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "sulabs.ml",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name sulabs.ml encrypt.ml;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
    #ssl_certificate /etc/letsencrypt/live/sulabs.ml/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/sulabs.ml/privkey.pem;
    root /srv/www/sulabs;
    location / {
        index index.html;
    }
    location ~* \.(js|css|json|png|jpg|jpeg|gif|ico)$ {
        expires 6h;
        add_header Cache-Control "public, no-transform";
    }
}
`,
		},
		{
			name: "test for hosting files without index with caching, HSTS and custom max age",
			exec: func() (string, string) {
				additions := Additions{
					AddHSTSConfig:     true,
					AddSecurityConfig: false,
					MakeDefaultServer: false,
					AddCachingConfig:  true,
					MaxCacheAge:       "1d",
				}
				service := Service{
					Selection:  2,
					Domains:    "encrypt.ml",
					Root:       "/srv/www/encrypt.ml",
					URL:        "",
					Port:       443,
					Additional: additions,
				}
				return prepareServiceFileContents(service)
			},
			expectedFileName: "encrypt.ml",
			expectedFileContents: `server {
    #listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name encrypt.ml;
    access_log off;
    error_log /dev/null crit;
    #ssl_protocols TLSv1.2 TLSv1.3;
    #ssl_certificate /etc/letsencrypt/live/encrypt.ml/fullchain.pem;
    #ssl_certificate_key /etc/letsencrypt/live/encrypt.ml/privkey.pem;
    #Send HSTS header
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";
    location / {
        root /srv/www/encrypt.ml;
    }
    location ~* \.(js|css|json|png|jpg|jpeg|gif|ico)$ {
        expires 1d;
        add_header Cache-Control "public, no-transform";
    }
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
