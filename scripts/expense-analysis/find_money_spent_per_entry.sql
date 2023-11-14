-- find how much you spent for all entries (excluding income) in a certain
-- currency.
SELECT 
    en.category, en.subcategory, SUM(ex.amount) * POW(10, ex.exponent) as total_expenses, COUNT(ex.id) as no_expenses  FROM expenses ex join entries en on en.id = ex.entry_id 
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category <> 'income' AND ex.currency = 'RON'
GROUP BY 
    en.category, en.subcategory
ORDER BY
    en.category, en.subcategory DESC;


-- find out how much money you spend in total on all expenses
SELECT SUM(total_expenses) FROM (
    SELECT 
        en.category, en.subcategory, SUM(ex.amount) * POW(10, ex.exponent) as total_expenses, COUNT(ex.id) as no_expenses  FROM expenses ex join entries en on en.id = ex.entry_id 
    WHERE 
        en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category <> 'income' AND ex.currency = 'RON'
    GROUP BY 
        en.category, en.subcategory
    ORDER BY
        en.category, en.subcategory DESC
)


-- find how much you spent for all categories and the
-- number of expenses in each, in a currency.
SELECT 
    en.category, SUM(ex.amount) * POW(10, ex.exponent) as total_expenses, COUNT(ex.id) as no_expenses  FROM expenses ex join entries en on en.id = ex.entry_id 
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category <> 'income' AND ex.currency = 'RON'
GROUP BY 
    en.category
ORDER BY
    total_expenses DESC;

-- find all expenses of a particular entry, in a currency.
SELECT 
    en.category, en.subcategory, ex.name, ex.amount, ex.exponent, ex.currency  FROM expenses ex join entries en on en.id = ex.entry_id 
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category == 'music' AND en.subcategory == 'interface' AND ex.currency = 'RON'
ORDER BY
    ex.currency, en.amount DESC;
