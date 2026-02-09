import { getRecommendedProblems } from "./api";
import { getRatingColor } from "./charts";
import { observeAndAnimate } from "./entrance-anim";
import problemHTML from "./templates/problem.html";

const PROBLEM_COUNT = 6;
const CONTEST_LOOKBACK = 5;
let solvedByContests;

const problemTemplate = loadProblemTemplate();
const problemContainer = document.getElementById("problem-container");

// Rating slider update logic
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

export function recommendProblems(signal) {
	if (solvedByContests === undefined)
		return

	try {
		const ratingRange = getSelectedRatingRange();
		getRecommendedProblems(PROBLEM_COUNT, ratingRange, solvedByContests, signal).then(problems => {
			clearProblemContainer();
			const elements = new Array();

			problems.forEach(resp => {
				const problemData = {
					name: resp.problem.name,
					id: resp.problem.contestId + resp.problem.index,
					rating: resp.problem.rating,
				};

				const element = displayProblem(problemData, resp.problem.tags);
				elements.push(element);
			});

			observeAndAnimate(elements);
		});
	} catch (err) {
		console.error(`Encountered problem recommending problems: ${err}`);
	}
}

export function updateSubmissions(submissions) {
	submissions = filterSubmissionsRecentContests(submissions, CONTEST_LOOKBACK);
	solvedByContests = findSolvedProblemsRecentContests(submissions);
}

export function setRatingRange(min, max) {
	updateRatingRange(false, min - MIN_RATING, max - MIN_RATING);
}

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

function filterSubmissionsRecentContests(submissions, amountOfContests) {
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

function findSolvedProblemsRecentContests(contests) {
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

// @param problemData Object with the name, id, and rating properties.
function displayProblem(problemData, tags) {
	const problem = document.importNode(problemTemplate.content, true);
	const problemElem = problem.firstElementChild;

	// Update data values.
	Object.keys(problemData).forEach(key => {
		const elem = problem.querySelector(`[data-field="${key}"]`);
		if (elem)
			elem.textContent = problemData[key];
	});

	// Update rating color.
	const ratingElem = problem.querySelector(`[data-field="rating"`);
	const ratingColor = getRatingColor(problemData.rating)
	ratingElem.style.setProperty("--text-color", ratingColor);

	// Add tags to the tags container.
	const tagsContainer = problem.querySelector(`[data-container="tags"]`);
	tags.forEach(tag => {
		const elem = document.createElement("span");
		elem.className = "problem-tag";
		elem.textContent = tag;
		tagsContainer.appendChild(elem);
	});

	problemContainer.appendChild(problem);
	return problemElem;
}

function getSelectedRatingRange() {
	return {
		min: parseInt(rangeInputs[0].value) + MIN_RATING,
		max: parseInt(rangeInputs[1].value) + MIN_RATING,
	};
}

function loadProblemTemplate() {
	const parser = new DOMParser();
	const doc = parser.parseFromString(problemHTML, "text/html");
	const template = doc.getElementById("problem-template");
	return template;
}

function clearProblemContainer() {
	while (problemContainer.firstChild)
		problemContainer.removeChild(problemContainer.lastChild);
}
