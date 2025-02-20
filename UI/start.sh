#!/bin/sh

# Replace environment variables in the Nginx configuration template
envsubst '${API_URL}' < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/default.conf

# Start Nginx
nginx -g 'daemon off;'
