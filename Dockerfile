FROM umputun/baseimage:buildgo as builder

ARG CI
ARG CI_COMMIT_BRANCH
ARG CI_MERGE_REQUEST_PROJECT_ID
ARG CI_COMMIT_SHORT_SHA

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

RUN apk --no-cache add ca-certificates
COPY . /src
WORKDIR /src

# install gcc in order to be able to go test package with -race
RUN apk --no-cache add gcc libc-dev

RUN \
    if [ -z "$CI" ] ; then echo "runs outside of CI" && version="$(/script/git-rev.sh)" ; \
    else version=${CI_COMMIT_BRANCH}${CI_MERGE_REQUEST_PROJECT_ID}-${CI_COMMIT_SHORT_SHA:0:7}-$(date +%Y%m%d-%H:%M:%S) ; fi && \
    echo "version=$version" && \
    go build -o /zenb -ldflags "-X main.revision=${version} -s -w" .

FROM scratch

COPY --from=builder /zenb /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/zenb"]
