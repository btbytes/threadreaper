# Thread Reaper

A simple utitlity to convert [Thread Reader](https://threadreaderapp.com/) pages into simpler HTML pages.

Usage:

	./threadreaper [-title TITLE] [-author AUTHOR] [-css CSSURL] URL  > local_page.html
	./threadreaper -title "On Binbary Numbers" -author foone https://threadreaderapp.com/thread/1202293011150852096.html  > foone_binary.html

This program also fixes an issue with embedded images that you would have if you were to simply download the Thread Reader HTML to disk.

If not provided,  the output file name will be decided based on the author name and the title.