#Goquadro Project description

[![Build Status](https://secure.travis-ci.org/goquadro/core.png)](http://travis-ci.org/goquadro/core)

[Goquadro][goquadro] is a web-based, workflow-aware cloud storage system.

The web interface to Goquadro is provided by a single-page application built with **AngularJS** alongside with **HTML5** and **CSS3**.  
The [secure connection][qualys-test] to the servers is powered by **nginx** in a **forward-security**-compliant **HSTS** mode. Goquadro servers are fully [IPv6 compliant](http://ipv6-test.com/validate.php?url=www.goquadro.com).  
Cookieless authentication is performed using **JSON web tokens** and persisted through HTML5 **web storage**.  
The api server is written in **Go** and serves data to the application through a **JSON** interface.  
User passwords are stored using a **bcrypt** implementation.  
Application data is stored in a **MongoDB** database, whereas user files are stored on Amazon S3.  
Goquadro sends emails (for account verification and for user notification) using the [Mailgun][mailgun] API.  
Goquadro sits on a dedicated infrastructure, managed via **ssh** and monitored using [NewRelic](http://newrelic.com).  
Basic interaction with Goquadro is provided through the dedicated [Google Chrome extension](http://bit.ly/gqchromeext), which lets the user push to Goquadro a reference to webpages or remote files with a single click.

All the code was written using SublimeText3 and VIM, versioned using **git** and saved on a GitLab instance installed on a virtual private server.


Credits
=======

Goquadro was built using opensource software.

* [Go][golang.org] by Google ([MIT license](http://golang.org/LICENSE))
* Javascript
* [nginx][nginx] by Igor Sysoev ([BSD license](http://nginx.org/LICENSE))
* [MongoDB][mongodb] ([GNU Affero General Public License](http://www.gnu.org/licenses/agpl-3.0.html))
* [mgo][mgo] MongoDB driver for Go by Gustavo Niemeyer ([Simplified BSD License](http://en.wikipedia.org/wiki/BSD_licenses#2-clause_license_.28.22Simplified_BSD_License.22_or_.22FreeBSD_License.22.29))
* [Martini][martini] by the Code Gangsta ([MIT license](https://raw.githubusercontent.com/go-martini/martini/master/LICENSE))
* [AngularJS][angularjs] by Google ([MIT license](https://raw.githubusercontent.com/angular/angular.js/master/LICENSE))
* [bcrypt implementation][go.bcrypt] by the Go authors


[goquadro]: https://www.goquadro.com "Goquadro.com"
[go-authors]: https://golang.org/AUTHORS "List of he Go authors"
[angularjs]: https://angularjs.org/ "AngularJS website"
[nginx]: http://nginx.org/ "nginx website"
[mongodb]: http://www.mongodb.org/ "mongoDB website"
[qualys-test]: https://www.ssllabs.com/ssltest/analyze.html?d=goquadro.com "Qualys SSL Labs report on https://www.goquadro.com"
[golang.org]: http://golang.org/ "Golang website"
[mgo]: http://labix.org/mgo "Rich MongoDB driver for Go"
[labix]: http://labix.org/ "Gustavo Niemeyer website"
[go.bcrypt]: http://code.google.com/p/go.crypto/bcrypt
[martini]: http://martini.codegangsta.io/ "Framework for web developing in Go"
[mailgun]: http://www.mailgun.com/ "Mailgun by Rackspace"
[codegangsta]: http://codegangsta.io/ "The Code Gangsta website"
[igorsysoev]: http://sysoev.ru/en/ "Igor Sysoev website"
