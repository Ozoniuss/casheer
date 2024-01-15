-- insert entry

-- all entries
SELECT * FROM entries WHERE month=12 AND year=2023 ORDER BY category, subcategory, amount DESC;

-- all entries, in RON, excluding income and economy
SELECT 
    en.id, en.category, en.subcategory, en.amount * POW(10, en.exponent) as expected_total, en.currency
FROM entries en
WHERE 
    en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now'))  AND en.category <> 'income' AND en.category <> 'economy' AND en.currency = 'RON'
ORDER BY 
    category, subcategory, amount DESC;

-- sum of all expected expenses, in RON, excluding economy and income completely.
SELECT SUM(total)
FROM (
    SELECT 
        en.category, SUM(en.amount * POW(10, en.exponent)) as total, en.currency
    FROM entries en
    WHERE 
        en."month" = strftime('%m', date('now')) AND en."year" = strftime('%Y', date('now'))  AND en.category <> 'income' AND en.category <> 'economy' AND en.currency = 'RON'
    GROUP BY 
        en.category
    ORDER BY 
        category, subcategory, amount DESC
);
