# Function: filter unique solutions
This function is for filtering out duplicate solutions and unsolved problems from a list of submissions.


```text
define function filterSolved:
	let solved be an empty list
	for every submission:
		if the submission was accepted:
			add new property "problem id" to the submission:
				contest id concatenated the problem index
			add the submission to solved

	sort the solved list based on the problem id property
	let uniqueSolved be an empty list
	for every submission in solved (in sorted order):
		if (uniqueSolved is empty) or (submission problem id ≠ last unique-solved problem id):
			add the submission to the end of uniqueSolved

	return uniqueSolved
```
