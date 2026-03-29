import { getRecommendedProblems } from "./api";
import { getRatingColor } from "./charts";
import { observeAndAnimate } from "./entrance-anim";
import { showError } from "./error";
import problemHTML from "./templates/problem.html";

const PROBLEM_COUNT = 6;
const CONTEST_LOOKBACK = 5;
let solvedSubs;

const problemTemplate = loadProblemTemplate();
export const problemContainer = document.getElementById("problem-container");

// Rating slider update logic
const MIN_RATING = 800;
const rangeSlider = document.getElementById("range-slider");
const rangeInputs = document.querySelectorAll(".range-input input");
const rangeTexts = document.querySelectorAll("#rating-range>input");

const STALE_OPACITY = 0.3;
const INPUT_WAIT = 500; // Time in ms.
let inputTimer;
let controller = new AbortController();

updateRatingRange(false, rangeInputs[0].value, rangeInputs[1].value);

rangeInputs.forEach(input => {
	input.addEventListener("input", e => {
		const rangeMin = rangeInputs[0].value;
		const rangeMax = rangeInputs[1].value;
		const isMin = e.target.dataset.rangeType === "min";

		updateRatingRange(isMin, rangeMin, rangeMax);

		clearTimeout(inputTimer);
		inputTimer = setTimeout(recommendProblems, INPUT_WAIT, solvedSubs);
	});
});

rangeTexts.forEach(input => {
	input.addEventListener("change", e => {
		const rangeMin = rangeTexts[0].value - MIN_RATING;
		const rangeMax = rangeTexts[1].value - MIN_RATING;
		const isMin = e.target.dataset.rangeType === "min";

		updateRatingRange(isMin, rangeMin, rangeMax);

		recommendProblems(solvedSubs);
	});
});

export function recommendProblems(solved) {
	if (solved === undefined)
		return;

	// Cancel previous requests.
	controller.abort();
	controller = new AbortController();

	solvedSubs = solved;

	setProblemOpacity(STALE_OPACITY);

	const ratingRange = getSelectedRatingRange();
	getRecommendedProblems(PROBLEM_COUNT, ratingRange, CONTEST_LOOKBACK, solved, controller.signal).then(problems => {
		const elements = new Array();

		const scrollY = window.scrollY;
		clearProblemContainer();

		problems.forEach(resp => {
			const element = displayProblem(resp.problem);
			elements.push(element);
		});

		observeAndAnimate(elements);

		window.scrollTo({ top: scrollY, behavior: "instant" });
	}).catch(err => {
		showError(err, problemContainer);
	});
}

export function setRatingRange(min, max) {
	updateRatingRange(false, min - MIN_RATING, max - MIN_RATING);
}

export function clearProblemContainer() {
	while (problemContainer.firstChild)
		problemContainer.removeChild(problemContainer.lastChild);
}

function updateRatingRange(isMin, rangeMin, rangeMax) {
	rangeMin = clampRating(rangeMin);
	rangeMax = clampRating(rangeMax);

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

function clampRating(rating) {
	rating = Math.max(rating, 0);
	rating = Math.min(rating, rangeInputs[0].max);
	return rating;
}

// @param problemData Object with the name, constestId, index, rating, and tags properties.
function displayProblem(problem) {
	const problemClone = document.importNode(problemTemplate.content, true);
	const problemElem = problemClone.firstElementChild;

	problem.id = problem.contestId + problem.index;

	// Update data values.
	Object.keys(problem).forEach(key => {
		const elem = problemClone.querySelector(`[data-field="${key}"]`);
		if (elem)
			elem.textContent = problem[key];
	});

	// Update rating color.
	const ratingElem = problemClone.querySelector(`[data-field="rating"`);
	const ratingColor = getRatingColor(problem.rating)
	ratingElem.style.setProperty("--text-color", ratingColor);

	// Add tags to the tags container.
	const tagsContainer = problemClone.querySelector(`[data-container="tags"]`);
	problem.tags.forEach(tag => {
		const elem = document.createElement("span");
		elem.className = "problem-tag";
		elem.textContent = tag;
		tagsContainer.appendChild(elem);
	});

	// Set anchor URL.
	const cfProblemsetURL = "https://codeforces.com/problemset/problem/";
	const anchorElem = problemClone.querySelector(`[data-link="problem-link"]`);
	anchorElem.href = cfProblemsetURL + problem.contestId + "/" + problem.index;

	problemContainer.appendChild(problemClone);
	return problemElem;
}

function setProblemOpacity(opacity) {
	for (let i = 0; i < problemContainer.children.length; i++) {
		const child = problemContainer.children[i];
		child.style.opacity = opacity;
	}
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
