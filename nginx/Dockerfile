FROM nginx:stable-bookworm

RUN apt-get update && \
	apt-get install mkcert && \
	rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/nginx/certs && \
	mkcert -cert-file /etc/nginx/certs/localhost.pem -key-file /etc/nginx/certs/localhost-key.pem \
	localhost 127.0.0.1

COPY nginx.conf /etc/nginx/conf.d/default.conf
