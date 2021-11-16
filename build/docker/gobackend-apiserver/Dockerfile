FROM BASE_IMAGE
LABEL maintainer="<i@sre.im>"

WORKDIR /app

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone 

COPY gobackend-apiserver /app/
COPY *.yaml /app/

ENTRYPOINT ["./gobackend-apiserver"]
CMD ["-c", "./gobackend-apiserver.yaml"]
