# gowiki
*A simple web-application written in Go, using a guide from golang.org*

**What is Gowiki?** 

Gowiki is a web-application that serves as a very simple wiki-site. 

**How does it work?**

When Gowiki is compiled, the binary can be executed from whatever folder you designated Go to store binaries in (Typically /bin). This then sets up a web-server that serves accessing users wiki-pages that they can edit, and interlink. 

**Requirements:**

To make this work, you must include the folders `data` and `tmpl` (Templates) within the directory from which you execute the binary. For instance, if the required folders are in `/Documents/Go/src/gowiki`, and your binary is placed inside `/Documents/Go/bin`, then you must CD to `/Documents/Go/src/gowiki`, and then run the binary by specifying its path (For instance, `$HOME/Documents/Go/bin/gowiki`, where gowiki is the name of the executable binary). 
