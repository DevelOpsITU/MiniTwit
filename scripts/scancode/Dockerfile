FROM python:3.9

RUN apt-get update && apt-get install bzip2 xz-utils zlib1g libxml2-dev libxslt1-dev

ADD "https://github.com/nexB/scancode-toolkit/releases/download/v30.1.0/scancode-toolkit-30.1.0_py39-linux.tar.xz" .

RUN mkdir scancode-toolkit && tar -xvf scancode-toolkit-* -C scancode-toolkit --strip-components=1

WORKDIR scancode-toolkit

RUN ./scancode --help

ENV PATH=$HOME/scancode-toolkit:$PATH

CMD ./scancode -clipeu --verbose /scan --processes `expr $(nproc --all) - 1` --json /scan/licenses.json

# From https://github.com/beevelop/docker-scancode/blob/master/Dockerfile