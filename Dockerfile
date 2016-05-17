FROM quay.io/brianredbeard/corebox

EXPOSE 8000
CMD ["/bin/reservation", "serve", "-a", "0.0.0.0:8000"]

ADD reservation /bin/reservation
