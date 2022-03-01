#!/bin/sh
docker run --name=product-review --network=product_review_default  -p 80:80 -it product-review:latest