
all: html

OPTIONS=--self-contained --mathml
#OPTIONS=--standalone --mathjax

RST2HTML=pandoc ${OPTIONS} --css=./css/pandoc.css -t slidy
RST2PDF=pandoc ${OPTIONS} --latex-engine=xelatex --css=./css/pandoc.css

talk01-srcfiles=\
	talk01.00.header.rst \
	talk01.01.userguide.rst \
	talk01.02.gettingstarted.rst \
	talk01.03.basicdesign.rst

talk01.html: ${talk01-srcfiles} css/pandoc.css
	$(RST2HTML) ${talk01-srcfiles} -o $@

view: talk01.html
	firefox $<

clean:
	rm -f *.html *.pdf
	find . -name "*~" | xargs rm -f

talk01.pdf: ${talk01-srcfiles}
	$(RST2PDF) ${talk01-srcfiles} -o $@

html: talk01.html
pdf: talk01.pdf
