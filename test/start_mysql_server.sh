#! /bin/sh

MYSQL_NAME=mysql
PORT=3306
MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=test
MYSQL_USER=test
MYSQL_PASSWORD=test

docker run -d --name $MYSQL_NAME --restart always \
    -h mysql-server \
    -p 3306:3306 \
    --env MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
    --env MYSQL_DATABASE=$MYSQL_DATABASE \
    --env MYSQL_USER=$MYSQL_USER \
    --env MYSQL_PASSWORD=$MYSQL_PASSWORD \
    mysql:5.7