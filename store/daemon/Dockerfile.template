#FROM is generated dynamically by the Makefile

EXPOSE 6800 6801 6802 6803

ADD bin/boot /app/bin/boot
ENTRYPOINT ["/app/bin/boot"]

# remove osd from copy-on-write
VOLUME /var/lib/ceph/osd

CMD ["ceph-osd"]
