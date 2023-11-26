FROM nginx:latest

COPY nginx.conf /etc/nginx/nginx.conf
COPY content/ /content/
