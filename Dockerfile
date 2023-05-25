FROM golang:1.20-alpine
COPY ./score /
EXPOSE 3000
CMD [ "/score" ]
