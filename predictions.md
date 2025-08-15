# Predictions explained

This file explains the typical entries I would write for expenses.

## Format

Each entry must have the `expected_total` field set. This represents the estimated value for that subcategory for the period. 

The other optional fields are:

- `recurring` tells casheer if this entry should be included in the next period's generated template
- `exponent` can be set to interpret the `expected_total` value differently (by default it's specified in the currency's major unit, having up to two decimals).
- `currency` can be set if that expense is managed in a different currency

## List of categories and subcategories

Entries have a category and subcategory attached to it to uniquely identify the entry. Category can be used for very high-level statistics and subcategory defines uniqueness. Every expense will fit in exactly one entry.

You may define your own categories depending on what you want to track. For example, you may want to track food very strictly as a separate category, or you're ok with food being added as a subcategory in a general shopping category.

Additionally, subcategories can be used to classify one-off purchases that aren't typical for a regular month. For example, for an `online` category represeting online shopping, one could add `tv` as a subcategory for that month.

Remember that these values are predictions for the following period, not actual expenses. Expenses will be introduced subsequently, either to one of these entries, or to a new one (e.g. a one-off purchase on the spot).

Here is an example list of categories that may be present in a regular template:

### Income

A special category storing the money coming in. Subcategories represent different sources of income, such as:

- `salary` for salaries employees
- `scolarship` for students
- `rent` for landlords

### Economy

This category is related to money management, primarily investments. Subcategories include:

- `etfs` for your monthly etf investments
- `bitcoins` if you invest in cryptos
- `donations` for a monthly donation
- `emergency_fund` if you're topping your emergency fund
- `extra_saving` e.g. if saving for a vacation


### Food

I like to track food more strictly to understand what I eat and what I spend on eating. This includes:

- `groceries` includes regular groceries (milk, water, yoghurt etc.)
- `snacks` tracks my addiction to snacks and juice
- `takeaway` represents how much I ordered online

I rarely ever eat out, so for example here `restaurant` could be a one-off category

### Bills

Represents my home bills.

- `digi` represents my phone and internet bill
- `electricity` is self explanatory
- `gas` is self explanatory
- `heat` is self explanatory
- `water` is self explanatory

### Transport

How much I spend for transport. My electric scooter helps keep this low.

- `fuel` is for my car fuel
- `parking` represents my monthly parking payments
- `rides` how much I spend on bus, taxi, ridesharing etc.


### Subscriptions

All my subscriptions. Helps me know when I get overboard.

- `youtube` probably my favourite subscription
- `backblaze` for my storage
- `google_drive` for my photos
- `medium` for my medium account
- `gergely` for "the pragmatic engineer"

### Fun

An umbrella coin for all my spending during fun events such as hangouts. Depending from context, expenses from other entries may end up here instead. For example, if I eat a pizza while I'm having a beer, I will not include it in `food`

- `going_out` general hanging out, typically beer
- `coffee` I regularly have coffee with my friends in the morning

### Health

I have ADHD which means I'll make regular health payments. But, I also don't really take good care of my health and this helps me build good habits.

- `adhd` my adhd meds
- `general` general pharmacy stuff

Other subcategories may include one-offs like a doctor's appoinment.

### Misc

Things that I'm not particularly interested to evaluate under a more specific broad category. May also be things like randomly buying a router.

- `idk` money that I spend but forgot on what

Other entries may be defined in an on-demand basis. The ones here represent my tracking interests at the point of writing.
