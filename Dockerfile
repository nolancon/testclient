FROM  clearlinux/os-core:latest 
COPY ./testclient /usr/local/bin/
CMD ["/usr/local/bin/testclient"]
