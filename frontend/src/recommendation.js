import { getRatingColor } from "./charts";

const rangeSlider = document.getElementById("range-slider");
const rangeInputs = document.querySelectorAll(".range-input input");
const rangeTexts = document.querySelectorAll("#range>input");

rangeInputs.forEach(input => {
	input.addEventListener("input", e => {
		const rangeMin = parseInt(rangeInputs[0].value);
		const rangeMax = parseInt(rangeInputs[1].value);

		if (rangeMax < rangeMin) {
			if (e.target.dataset.rangeType === "min")
				e.target.value = rangeMax;
			else
				e.target.value = rangeMin;
		} else {
			rangeSlider.style.setProperty("--left-pos",
				(rangeMin / rangeInputs[0].max) * 100 + '%');
			rangeSlider.style.setProperty("--right-pos",
				((rangeInputs[0].max - rangeMax) / rangeInputs[0].max) * 100 + '%');
		}
	});
});

export function filterSubmissionsRecentContests(submissions, amountOfContests) {
	const recentContestsId = [];
	const recentContests = [];

	for (let i = 0; i < submissions.length; i++) {
		if (submissions[i].author.participantType != "CONTESTANT") {
			continue;
		}

		let ID = submissions[i].contestId;
		if (recentContestsId.includes(ID)) {
			recentContests[recentContestsId.indexOf(ID)].submissions.push(submissions[i]);
		} else if (recentContestsId.length < amountOfContests) {
			recentContestsId.push(ID);
			recentContests.push({ id: ID, submissions: [submissions[i]] })
		}
	}

	return recentContests
}

export function findSolvedProblemsRecentContests(contests) {
	const solvedByContests = [];

	for (const contest of contests) {
		const solvedInContest = [];

		for (const problem of contest.submissions) {
			if (problem.verdict == "OK" && !(solvedInContest.includes(problem.index))) {
				solvedInContest.push(problem.problem.index);
			}
		}

		solvedByContests.push({ id: contest.id, indices: solvedInContest })
	}

	return solvedByContests
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
