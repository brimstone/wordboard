FROM scratch

COPY app /

COPY static /

EXPOSE 8080

ENTRYPOINT ["/app"]
