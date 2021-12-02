FROM BASE_IMAGE
LABEL maintainer="<i@sre.im>"

WORKDIR /app

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone 

COPY gobackend-apiserver /app/
COPY *.yaml /app/

# Run in development/test/production by env value 'dev/test/prod' or 
# by specifing config file 'dev/test/prod.gobackend-apiserver.yaml'.
# E.g.: 
#    -e RUN_MODE=prod
#    or
#    -c prod.gobackend-apiserver.yaml
ENTRYPOINT ["./gobackend-apiserver"]
