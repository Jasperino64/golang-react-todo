FROM nginx:alpine

COPY dist /usr/share/nginx/html

# Make nginx listen on 8080 (Fly requirement)
RUN sed -i 's/listen\s\+80;/listen 8080;/' /etc/nginx/conf.d/default.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
