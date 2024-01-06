# UI

A very small I'm building for casheer. Not to be confused with the [real UI](https://github.com/iuliailies/casheer-client) my girlfriend is building (well, currenlty not building cause she's busy, that's why I started this one myself).

The need for this UI came after using the scripts I've talked about in my blog post (will update with link) from [this folder](../scripts/fill-expenses/). It's become more and more of a pain in the ass to use them. It works reasonably when I go through the expenses for the whole month and takes be between one and two hours, but my current goal is to spend 3 minutes each day adding everything new. And for that, those scripts are quite garbage. They're very annoying to work with, but the biggest problem is that there's no way at the moment to see my current progress.

I first wanted to do a visualization in the terminal, but that seemed pretty complicated and while I'd be really interested to see how it's done (such as in tools like [k9s](https://github.com/derailed/k9s)) it requires a significant time investment and it's not really helping me profesionally. Learning a bit of frontend however, is (and not the CLI "frontend"). I randomly stumbled upon a video using Go and HTMX together, particularly [this one](https://www.youtube.com/watch?v=F9H6vYelYyU), and I immediately decided to just go for it. It seemed at the moment just what I needed to hack this quickly in a few weekends.

I've never done frontend before, so I'm really curious to see how this turns out. Particularly, I'll be on the lookout for how fast I can pick it up (speed was one of the main reasons I decided to go with this approach), how much I get to learn about frontend and htmx and how much impact it brings to my day-to-day life. After all, if I won't be using this tool until I have a real frontend, I'm not getting that much value out of this project.

## How it's built

This is build using Go (particularly the html templating library) and HTMX. The idea is as follows:

- have the main API run on my machine;
- have another Go server running on my machine, serving the [templated HTML](./index.html);
- work in the browser, which sends the changes to the Go server, further sending the changes back to the main API.

See more about how [htmx works] and how [Go templating works](https://pkg.go.dev/text/template). Note that templating is a general concept and not specific to Go.

It does sound a bit weird to me that I'm having two servers just to create a database wrapper after all. At first glance I thought I'm doing something wrong because of this and there's too much work. But then I realised I'm actually doing _less_ work than I'd do by simply hacking a typical frontend that communicates with the backend directly.

And the primary reason for this is because I don't have a casheer domain type system implemented in any other language than Go. I can't communicate from javascript directly with Casheer unless I define [all the shitload](../pkg/casheerapi/) of types I have for the request/response dance, as well as the response handling logic. In Go, all my types already exist perfectly taylored to my API, and I crafted [a client](../client/) which handles all that logic. This client is also driving my [end-to-end](../e2e/) tests, so not only the client is tested, but it also helps test the API.

With my current approach, I just have to define some structs for the data I want to display on my page, which I would've needed to do anyway in javascript (or at least something similar). But in my bussiness logic, I can use the Go client. Since I'm using an HTML template, I'm making use of all that already existing Go logic to fill in the data-layer structs. So I just have my primary goal in mind: understanding how this frontend and templating works in general, and figuring out how to make the UI look. I don't need to think at all about the logic of interacting with the API.

## Other approaches?

Just for the record, I am unaware of any good library which can take code from a language and convert it to another, and I also didn't want to use swagger codegen since I already carefully hand-crafted all of my API-level data representation. If I remember correctly, I didn't have a good time with swagger codegen the limited number of times I tried it. Having such tool would bring multiple options to the table in terms of efficiency, but since the reseraching effort seemed bigger and with less learnings than the HTMX and templating route (and also less valuable since I already hand-crafted all my API-layer model representation), I took it as a waste of time.