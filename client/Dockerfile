# FROM node:20.12.0 AS base
# LABEL MAINTAINER="ivan.xiewei@foxmail.com"

# # 设置内存限制
# ENV NODE_OPTIONS="--max-old-space-size=4096"

# WORKDIR /app
# COPY . .

# RUN yarn install --registry=https://registry.npmmirror.com

# RUN yarn run build:prod

FROM nginx:alpine
LABEL MAINTAINER="ivan.xiewei@foxmail.com"
COPY ./.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/my.conf
COPY dist /usr/share/nginx/html
RUN ls -al /usr/share/nginx/html