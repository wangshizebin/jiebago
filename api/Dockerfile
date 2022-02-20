FROM alpine
ENV TZ=Asia/Shanghai
RUN mkdir -p /jiebago/dictionary
COPY ./jiebago /jiebago
COPY ./dictionary /jiebago/dictionary
WORKDIR /jiebago
CMD ["./jiebago -httpaddr=:8118"]