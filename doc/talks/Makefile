
OPTIONS=--self-contained --mathml
#OPTIONS=--standalone --mathjax

RST2HTML=pandoc ${OPTIONS} --css=./css/pandoc.css -t slidy

talk01.html: talk01.rst css/pandoc.css
	$(RST2HTML) $< -o $@

view: talk01.html
	firefox $<

clean:
	rm -f *.html
	find . -name "*~" | xargs rm -f
