FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf
ADD ./default.conf /etc/nginx/conf.d/default.conf

RUN mkdir /certs
# ADD ./fullchain.pem /certs/fullchain.pem
# ADD ./privkey.pem /certs/privkey.pem
