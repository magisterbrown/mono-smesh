FROM node:20 as build
COPY ./app /app
WORKDIR /app
RUN npm install
RUN npm run build

FROM nginx:latest
COPY --from=build /app/.svelte-kit/output/ /usr/share/nginx/leadapp
COPY app/favicon.ico /usr/share/nginx/leadapp/prerendered/pages
COPY nginx.conf /etc/nginx/conf.d/
#CMD nginx -c /etc/nginx/nginx.conf
CMD ["nginx", "-g", "daemon off;"]
