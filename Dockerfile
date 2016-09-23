FROM scratch

COPY app /

COPY static /static

EXPOSE 8080

ENTRYPOINT ["/app"]
