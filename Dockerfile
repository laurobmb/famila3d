FROM alpine:latest
WORKDIR /app
COPY dist/familia3d.bin ./familia3d.bin
COPY templates/ ./templates/
COPY static/ ./static/
RUN chmod +x ./familia3d.bin
EXPOSE 8080
CMD ["./familia3d.bin"]