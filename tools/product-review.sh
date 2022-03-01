#!/bin/sh
docker run --name=product-review --network=product-review_default -e PR_USERNAME="$PR_USERNAME" -e PR_PASSWORD="$PR_PASSWORD" -p 80:80 -it product-review:latest