FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY ./dist ./app/
WORKDIR /app
CMD ["./dailydo"]