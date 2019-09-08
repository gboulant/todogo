
all: html

#OPTIONS=--self-contained --mathml
OPTIONS=--standalone --mathjax

RST2HTML=pandoc ${OPTIONS} --css=./css/pandoc.css -t slidy
RST2PDF=pandoc ${OPTIONS} --latex-engine=xelatex --css=./css/pandoc.css

talk01.html: talk01.rst css/pandoc.css
	$(RST2HTML) $< -o $@

view: talk01.html
	firefox $<

clean:
	rm -f *.html *.pdf
	find . -name "*~" | xargs rm -f

talk01.pdf: talk01.rst
	$(RST2PDF) $< -o $@

html: talk01.html
pdf: talk01.pdf
