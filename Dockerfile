FROM almir/webhook:latest
# TODO: dont use latest tag

RUN apk --no-cache add python3 py-requests

ADD . /webhooks

# ENV GT_API_BASE
# ENV BK_TOKEN
# ENV GT_TOKEN
# ENV GT_SECRET
# ENV BK_SECRET

CMD ["-verbose", "-template", "-hooks=/webhooks/hooks.json"]
