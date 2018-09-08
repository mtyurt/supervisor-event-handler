FROM alpine

RUN set -xe && \
    apk add supervisor inotify-tools && \
	rm -rf /var/cache/apk/* /tmp/* && \
    mkdir -p /var/log/supervisor

ENTRYPOINT ["supervisord"]
CMD ["--nodaemon", "--configuration", "/root/eventhandler/conf/supervisord.conf"]
