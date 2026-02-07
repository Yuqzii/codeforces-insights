import { getRatingColor } from "./charts";

const MIN_RATING = 800;
const rangeSlider = document.getElementById("range-slider");
const rangeInputs = document.querySelectorAll(".range-input input");
const rangeTexts = document.querySelectorAll("#rating-range>input");

updateRatingRange(false, rangeInputs[0].value, rangeInputs[1].value);

rangeInputs.forEach(input => {
	input.addEventListener("input", e => {
		const rangeMin = rangeInputs[0].value;
		const rangeMax = rangeInputs[1].value;
		const isMin = e.target.dataset.rangeType === "min";

		updateRatingRange(isMin, rangeMin, rangeMax);
	});
});

rangeTexts.forEach(input => {
	input.addEventListener("change", e => {
		const rangeMin = rangeTexts[0].value - MIN_RATING;
		const rangeMax = rangeTexts[1].value - MIN_RATING;
		const isMin = e.target.dataset.rangeType === "min";

		updateRatingRange(isMin, rangeMin, rangeMax);
	});
});

function updateRatingRange(isMin, rangeMin, rangeMax) {
	rangeMin = Math.max(rangeMin, 0);
	rangeMax = Math.min(rangeMax, rangeInputs[0].max);

	if (rangeMax < rangeMin) {
		if (isMin)
			rangeMin = rangeMax;
		else
			rangeMax = rangeMin;
	}

	rangeTexts[0].value = rangeMin + MIN_RATING;
	rangeTexts[1].value = rangeMax + MIN_RATING;

	rangeInputs[0].value = rangeMin;
	rangeInputs[1].value = rangeMax;

	// Update visual slider
	rangeSlider.style.setProperty("--left-pos",
		(rangeMin / rangeInputs[0].max) * 100 + '%');
	rangeSlider.style.setProperty("--right-pos",
		((rangeInputs[0].max - rangeMax) / rangeInputs[0].max) * 100 + '%');

	// Update number input colors
	rangeTexts[0].style.color = getRatingColor(rangeMin + MIN_RATING);
	rangeTexts[1].style.color = getRatingColor(rangeMax + MIN_RATING);
}

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
