import { getRatingColor } from "./charts";

export function filterSubmissionsRecentContests(submissions) {
	const filtered = [];
	const recentContests = [];

	for (let i = 0; i < submissions.length; i++) {
		if (submissions[i].author.participantType != "CONTESTANT") {
			continue;
		}

		let amountOfContests = 5;
		let ID = submissions[i].contestId;
		if (recentContests.includes(ID)) {
			filtered[recentContests.indexOf(ID)].push(submissions[i]);
		} else if (recentContests.length < amountOfContests) {
			filtered.push([submissions[i]]);
			recentContests.push(ID);
		}

	}

	return filtered
}

// Sets the color of all problem rating elements to their corresponding Codeforces rank color.
export function updateProblemRatingColors() {
	const elements = document.querySelectorAll(".problem-rating");
	elements.forEach(element => {
		const rating = parseInt(element.innerText);
		const color = getRatingColor(rating);
		element.style.setProperty("--text-color", color);
	});
}
