FROM alpine:3.9

RUN sed -i -e 's/v3\.4/edge/g' /etc/apk/repositories \
    && apk upgrade --update-cache --available \
    && apk --no-cache add librdkafka

WORKDIR /app/
COPY ./bin/supermart /app/
RUN ls

RUN chmod +x supermart

CMD["./supermart","--host=0.0.0.0", "--port=5000"]
