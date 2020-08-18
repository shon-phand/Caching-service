# Caching-service

In this project Golang is used to for creating rest endpoints. Redis is used for caching the data. Data is saved in mysql database.

To run the projrct please run docker-compose.yml

$docker-compose up -d

You would need to run first /internal/currency/UpdateDb api endpoint to update the DB(mysql) to fill the data in mysql database.

There are 2 endpoints available to users

    /currency/:symbol : to get the single currency info.
    /currency : to get the all currency with info.

Initially cache will be empty, once the request made to paticular currecy, currency will be fetched from the mysql database and will be cached for 500 seconds. If user trigger all currency then all currencies will be store in cache for 500 seconds.

If within 500 seconds another request made for the data in cache data will be served from cache.
