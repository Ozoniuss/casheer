-- all income
SELECT
	subcategory, amount * POWER(10, exponent), currency, recurring  
FROM
	entries
WHERE
	category = 'income';

-- income during current month
SELECT
	subcategory, amount * POWER(10, exponent), currency, recurring
FROM
	entries
WHERE
	category = 'income' AND "month"=10 AND "year"=2023