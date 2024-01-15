-- all incomes
SELECT
	subcategory, amount * POWER(10, exponent), currency, recurring  
FROM
	entries
WHERE
	category = 'income';

-- predicted income during current month
SELECT
	subcategory, amount * POWER(10, exponent), currency, recurring
FROM
	entries
WHERE
	category = 'income' AND "month"=strftime('%m', date('now')) AND "year"=strftime('%Y', date('now')) 

-- predicted vs actual income current month
SELECT
	en.category, en.subcategory, en.amount * POWER(10, en.exponent), en.currency, SUM(ex.amount) * POW(10, ex.exponent) as total_received
FROM
	entries en LEFT JOIN expenses ex on ex.entry_id = en.id
WHERE
	en.category = 'income' AND en."month"=strftime('%m', date('now')) AND en."year"=strftime('%Y', date('now')) 
GROUP BY
	en.category, en.subcategory;