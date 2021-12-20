From golang:1.17.4
COPY main main
EXPOSE 3000
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["./main"]
