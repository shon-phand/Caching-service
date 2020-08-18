FROM mysql

ENV MYSQL_DATABASE currency_db
ENV MYSQL_ROOT_PASSWORD mysqlpass
COPY ./sql-scripts/ /docker-entrypoint-initdb.d/

ADD ./sql-scripts/ /docker-entrypoint-initdb.d/

