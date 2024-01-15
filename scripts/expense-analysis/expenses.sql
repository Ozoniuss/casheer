-- all current expenses for all entries, excluding income and economy, in RON
SELECT 
    en.category, en.subcategory, en.amount * POW(10, en.exponent) as expected_total, SUM(ex.amount) * POW(10, ex.exponent) as total_expenses, COUNT(ex.id) as no_expenses  FROM entries en LEFT JOIN expenses ex on en.id = ex.entry_id 
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category <> 'income' AND en.category <> 'economy' AND ex.currency = 'RON'
GROUP BY 
    en.category, en.subcategory
ORDER BY
    en.category, en.subcategory DESC;


-- running total for all entries, excluding income and economy, in RON
-- this is just the sum of the above
SELECT SUM(total_expenses) FROM (
    SELECT 
        en.category, en.subcategory, SUM(ex.amount) * POW(10, ex.exponent) as total_expenses, COUNT(ex.id) as no_expenses  FROM expenses ex join entries en on en.id = ex.entry_id 
    WHERE 
        en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category <> 'income' AND en.category <> 'economy' AND ex.currency = 'RON'
    GROUP BY 
        en.category, en.subcategory
    ORDER BY
        en.category, en.subcategory DESC
);


-- find all expenses of a particular entry, in a currency.
SELECT 
    en.category, en.subcategory, ex.name, ex.amount * POW(10, ex.exponent) as amount, ex.exponent, ex.currency  FROM expenses ex join entries en on en.id = ex.entry_id 
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now')) AND en.category == 'fun' AND en.subcategory == 'bucharest' AND ex.currency = 'RON'
ORDER BY
    ex.currency, en.amount DESC;
