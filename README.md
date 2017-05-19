# HttpGrabber
HTTP grabber for capturing resource data from a single web URL

The motivation behind the grabber is to capture JSON data when ever a well known web page changes its content.

Output from the grabber will be written into a directory named using the seconds since the tool has been initiated.


## Build
To build use the golang software distribution for your hardware platform then:

<pre>
go get -d . && go build -o bin/HttpGrabber .
</pre>
