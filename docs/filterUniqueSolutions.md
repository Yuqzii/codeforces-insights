```text
define function filterSolved:
	let solved be an empty list
  for every submission:
		if the submission was accepted:
			add new property "problem id" to the submission consisting of contest id and the problem id within the contest
			add the submission to solved

	sor the solved list based on the problem id property
	let uniqueSolved be an empty list
	for every submission in solved (in sorted order):
		if (uniqueSolved is empty) or (this submissions problem id is unequal to the problem id of the last element of unique solved)
			add the submission to the end of uniqueSolved

  return uniqueSolved
```
