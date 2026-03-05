FROM gomicro/goose:3.26.0

COPY migrations/ /migrations/

RUN chmod -R 755 /migrations

ENTRYPOINT ["goose"]
