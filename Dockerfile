FROM public.ecr.aws/docker/library/golang:1.20.4-alpine
COPY ./score /
EXPOSE 3000
CMD [ "/score" ]
