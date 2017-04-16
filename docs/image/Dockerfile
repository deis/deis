FROM ubuntu-debootstrap:14.04

RUN apt-get update && apt-get install --yes python make python-dev libffi-dev git curl && apt-get clean

RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

ADD ../docs_requirements.txt /
ADD ./rundocs /

RUN pip install -r /docs_requirements.txt

WORKDIR /

CMD ["/rundocs"]
