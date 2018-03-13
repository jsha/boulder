FROM buildpack-deps:stretch-scm

ADD build.sh /tmp/build.sh
RUN mkdir go1.10
RUN curl -s https://dl.google.com/go/go1.10.linux-amd64.tar.gz | tar --strip-components 1 -C go1.10 -xpz

RUN mkdir go1.9
RUN curl -s https://dl.google.com/go/go1.9.linux-amd64.tar.gz | tar --strip-components 1 -C go1.9 -xpz

# Default Go version.
ENV PATH="/go1.9/bin:${PATH}"

RUN /tmp/build.sh
