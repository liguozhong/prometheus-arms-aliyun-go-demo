FROM BASEIMAGE

COPY version /
COPY ../cmd/entry /

RUN ls -hl /

#USER 1001:1001
ENTRYPOINT ["/entry"]
#CMD while true;do echo hello arms adapter;sleep 5;done
