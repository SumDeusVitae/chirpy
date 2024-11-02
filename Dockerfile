# Selecting lightweight Debian Linux OS
FROM debian:stretch-slim

# COPY source destination
COPY chirpy /bin/chirpy

# Set PORT environment variable
ENV PORT=8080

#Automatically start serer in the container when we run it
CMD ["/bin/chirpy"]
