FROM alpine

ADD terraform-gendocs /usr/local/bin
ADD output.tmpl /usr/local/bin
