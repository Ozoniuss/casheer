# The abstractions behind Casheer

Casheer provides multiple models (or resources, if you're a REST fanatic), but only two are fundamental to the core design of the application, those being "entries" and "expenses". There are also "debts" and I might introduce new models such as "wishes", but these are very easy to understand and it's not worth discussing them in a design document. At first glance, the former two might throw you off and it may not be that intuitive why they both exist, so I'll try my best to explain the reason behind them.

Thinking about it from a personal perspective
---------------------------------------------

I want to first address the issue I'm trying to solve with this application. I was never one to organize things, and finances is no different. A big part of my life, the only planning I did in terms of my financials was to look into my account and see how much money I have left. Later, I started investing and I quickly realised I needed more predicability. The jobs I was working back then were paying by the hour and thus didn't offer me a consistent monthly income, which made it difficult to plan consistent investments.

You may ask yourself: why not wait until the end of the month, and invest the money that's left? This sounds easy enough, but then:

- it forces me to wait until the end of the month, with the potential of missing good investment opportunities;
- it increases the possibility of reckless spending, knowing that you still have money left in your bank account.

I wanted to have a mechanism that will allow me to predict how much I could invest each month, without having to wait until the end of the month. Besides, I felt that being organized with my money would also encourage me to stick to my investments plan, which can get quite diverse (emergency fund, stocks, etfs, bank deposits, cryptocurrencies etc.). Another benefit of this approach is that it encourages me to stick to rational decisions: instead of doing emotional decisions during a month (e.g. crypto price is very low, let me buy as much as I can now) which I know I'm prone to, I would feel compelled to stick to my original plan and my long-term strategy. Of course, that is not to say I wouldn't get off track now and then when I notice great opportunities.

Besides, I know from past experiences how easy it is to let go of investments, even just for a month, and how difficult it is after that to start again. Having a good way to plan and track my spendings was going to build the discipline I lacked to start my path towards financial independence. I would never end up spending all the money I had earned that month just because it was easy to do so.

Now, all that was left for me was to come up with a reliable mechanism to accurately plan ahead what I was going to spend, to not go over my budget. I knew I wanted to plan and track in the timeframe of one month, mainly because I was paid monthly and I was dependent on my income stream. The simplest thing that came to mind was to just note down the total amount I spent each month, average it with the previous months and that was it. But I quickly realised this wasn't the best indicator, since the variance was quite high. Even the income of my second job was highly volatile. I needed something more granular that would allow me to make more accurate and consistent predictions, so I realised I need to group the expenses somehow. Then, the complexity shifts towards finding the right groups that aren't as volatile, which I felt much more comofortable with. Once I settled on this approach, those groups ended up what's being called "an entry".

The main idea is that for each month, you would define the "entries" that are most relevant to you, which are simply predictions on how much money would end up being spent on things that are related to that entry. There are very stable ones such as paying rent, buying groceries or fuel, buying etfs, which you can reliably use to make the predictions for the next month. Of course, there are also ones that aren't as stable, like one-time events such as going to concerts, going on holiday, or making an expensive purchase. These ones are typically not recurring and you can't really use them for the following month, but you also typically know about them beforehand at the beginning of the month, so they can be planned.

I was quite happy with this system. The only thing that was left was a mechanism to understand if my predictions was accurate, and for that I decided to simply store all purchases I was going to make during a month. Then, I was going to associate each purchase with an entry, and I could sum up all purchases for one entry to find out the actual amount of money I spent, compared to my original prediction. My setup was now complete, and I was able to get started on building my application.

I must confess, I spent quite a good deal of time choosing these two abstractions. There were many ways I could go about this problem, and at times it felt overwhelming. But in the end, I'm glad I eventually settled with those, as they provide me an enormous amount of flexibility. This is kindof how I was doing it at first on an Excel sheet, but it was definitely just intuition and I hadn't thought about it thoroughly. I can definitely understand Roy Fielding's [quote](https://roy.gbiv.com/untangled/2008/rest-apis-must-be-hypertext-driven) better now:

"However, I think most people just make the mistake that it should be simple to design simple things. In reality, the effort required to design something is inversely proportional to the simplicity of the result. As architectural styles go, REST is very simple."

Let's explore these two abstractions in depth, and understand their semantics in the context of Casheer.

The abstractions I chose
------------------------

The defining characteristics for an entry are a month, year, category, subcategory and value. In fact, the first 4 attributes have a UNIQUE constraint in the database. When it comes to value, for now you can think if it as just a number representing how much you're planning to spend. Though, I should mention defining values is also an extremely complicated topic which deserves its entire [design section.](#defining-value-is-hard)

Previously, I said that an entry is a prediction for a month. From this abstraction's perspective, when I say "month", I actually refer to "month" and "year" together, which is the basic unit of time measurement in Casheer. When you store an entry, it's referencing an actual (approximately) 30 consecutive day period, which is loosely the same period of a month in a certain year. I say "loosely" because you don't necessarily have to think of an entry being defined as the days of an exact month; that's just a convenience for me because I happen to make my predictions monthly.

As a side note, this abstraction can be easily generalized to simply refer to an interval of successive days, if someone prefers to organize their expenses over different time periods (e.g. weekly, or quarterly). But that was overkill for me, and the application assumes that you'd be doing it monthly by having a "month" column for the entries table, and validating that the month is between 1 and 12. This is a fairly easy fix if I ever plan to change it.

Note that it's not necessary to store the month and year as part of the entry directly. I've thought several times whether or not a better design would be to store periods separately, and associate entries to them, since it's very natural to think in terms of "the current month". In fact, my girlfriend who's doing the UI had a hard time to think of an entry this way, since you'd often (though not always) be creating an entry from the perspective of the current month and year, which will be part of the entry's data. But in the end I found this design to be very elegant, because illustrates what an entry tries to encapsulate perfectly, in a single abstraction: a prediction on how much you'd spend over a fixed period of time. The database doesn't care about how this correlates to a real month and year.

The last thing to talk about is category and subcategory. This probably creates the most confusion, so I'll try to do my best. Semantically, an entry in the context of casheer is a particular set of expenses that belong together, and that is meaningful to estimate. It is essentially a container that groups related expenses together, which added together give the actual amount of money spent compared to the original prediction. It is one unit of spending I'm particularly interested in, for a certain period

The entire entry is the "group" of expenses I was talking about in the first chapter. Not its category or subcategory alone, but both combined. The category is simply a high-level label that allows me to group expenses together into a bigger area, which combined with the subcategory helps me identify an area I'm particularly interested in predicting and tracking the spendings of the current month. I know this may seem overly-complicated, in which case I recommend reading the example below. But before that I want to touch on expenses.

An expense in the context of casheer stores an actual real transaction, and any details you want about it. Based on the period, it is associated with an entry, and we say it "contributes" to the running total of the entry. This is perhaps also a bit confusing at first: an entry can have any number of expenses. It can have one expense, if the entry represented a big purchase, or multiple smaller expenses, if I just wanted to track down something where I'm not interested in the individual transactions, but only the final outcome. Keep in mind that the entries are the abstractions that I'm actually interested in tracking, not expenses. Expenses help for "debugging" purposes, such as showing if the prediction is reliable, or breaking down a suspicious amount of money spent for an entry in smaller units, to understand what happened.

Let me give a few examples, to showcase how I intended for these abstractions to be used.

Real examples of using these abstractions
-----------------------------------------

Let's say that I am a smoker, and I want to make a prediction for how much money I will spend on cigarettes this month. I'm not interested in the individual transactions, I'm interested how much I will spend in total, so I can better understand the amount of money I'm left with to make some investments. It's October 2023, and I predict that I will buy cigarettes worth 100 EUR. Then, I record every individual cigarette purchase as an expense, which I associate with this particular entry.

At the end of the month, I can compare my prediction to the actual amount of money spent that comes from adding up all expenses for this entry. If they're very close, I know I have an accurate prediction. If it's significantly higher, I can go through all my expenses to understand how many packs I bought, and take action: either adjust the prediction for next month, or try to smoke less. The action is of course personal preference.

The category and subcategory is also specific to each person. If you're really passionate about smoking and are maybe an active collector, you can make several related entries that all have the main "smoking" category, but different subcategories: "cigarettes", "lighters", "ashtrays", "filters", "rolling papers" etc. This is helpful if you're interested in how much money you spend for smoking-related things as a whole, but you want to track down with more granularity.

However, if you're not really passionate about smoking and only buy cigarettes, your categories would look totally different. Maybe you'd have a single entry with a "shopping" category and "cigarettes" subcategory in this case, since you're only really interested in how much money goes away to smoking. It's important to understand when defining your entries that they should be specific to your use case, depending on what you want to track down and plan for.

I believe this example illustrates extremely well that even though the abstraction provides only a category and subcategory, it's not restrictive at all, which might be the general impression when thinking of online shops that have a lot of categories. The first example shows how you can have "smoking" as something very specific, which you could even further break down if desired into entries with categories such as "tobacco", and the second ones shows how you can think of smoking as something very general, that is compacted together with other items in a broader category. Unlike online shops which need to cater to a lot of people, this works nicely because the entries are specific to individuals. In my opinion, there is no need for extra levels of nesting online shops have.

Let's see how entries are used in a different context. Let's assume this month I will be buying a phone. Even if it's a one-time purchase, I will make a new entry, with either the exact price of the phone if I know it, or an estimate if maybe I'm buying second hand and I don't know the price beforehand. It's probably only going to have an expense associated to it, which will give me the exact price for accurate finance tracking. The category in this case is really user-specific, and could be anything from "tech" to "one-time purchases". Even if it's a single purchase, it's still something I can predict and use to help me plan my financials. If for example I make all my predictions and realise I'm over the budget, I will know to look for something cheaper. Or the opposite, I know I can increase my budget if I have some left-over money.

While these two examples point out the main use cases of entries, I want to give a third one that will illustrate what in my opinion is an "incorrect" way to use these abstractions. Let's say I'm tracking down all my uber trips in an entry with category "transport" and subcategory "uber", because I'm taking uber to work every day. But, this month I'm going to a concert which I'm also tracking. If going to the concert and back involves two uber trips, those trips should be recorded as expenses for that concert's entry, alongside tickets, maybe drinks etc. The reason I'm having the "trasnport, uber" entry is because I'm using uber regularly enough for it to make sense to track it as an individual entry. However, these two uber trips do not reflect the monthly reality of my regular uber use, but rather the one-time event, which only makes for a compelling argument that I should have the uber expenses as part of the concert entry. This will not affect my regular "transport, uber" prediction, and not only that, but tracking these trips and what I mentioned previously for the concert entry will give me a better estimate for the following concert.

Income is a category for an entry
---------------------------------

Perhaps surprisingly, I use "income" as a category for an entry. This is by design, and definitely reflects my early days of working when I wasn't having consistent income. Back then, it made sense to have predictions for my salary, because it wasn't consistent. As I illustrated in the previous examples, even if I have a consistent salary now, I still make entries for my income and I think of them in terms of predictions, and I also record expenses for them when I get paid. It's simply a more general abstraction that has wider use-cases, and is helpful both when the income is volatile, and also when there are multiple streams of income. In that case, you would have multiple entries, all with the "income" category.

The application itself doesn't differentiate between entries with the "income" category, and other entries. That is by design. Remember what I said about entries: they are simply a prediction I'm interested in and I want to track. From that point of view, any entry I make for income fits this purpose. This is why I decided to leave the clients interpret the meaning of "income" (i.e. exclude "income" entries if I want to calculate the total amount I'll be spending). In fact, it doesn't even need to be named income. It's each user's and each client's responsibility to interpret their sheet, not the API's, which just provides a means of storing these abstractions.

The goal is to be on zero
-------------------------

Everybody heard the saying "Rich people don't work for money", from the famous "Rich dad, poor dad" book. Money, like other things, is something of value. But on its own, it's something that typically depreciates in value, since supply increases. One of the goals when creating entries is to not be left with extra money that was not spent purposefully. Be it for your personal growth, investments, emergency fund, donations or whatever else, my goal is to have a balance equal to 0, once I deduct the value of all entries I created from my income. If it's not 0, it means that I have money lying around that I haven't used for anything purposeful to me that's just depreciating in value.

Defining value is hard
----------------------

I just want to share some of my thoughts when I was thinking about storing values. Of course, I'm no expert and I haven't studied economics, these being just my personal opinions. Until now, I avoided talking about values, and all I said is to think as it of how much money you're spending. You probably thought of it in your local currency.

TBD, as it's complicated and I'm still thinking about it.