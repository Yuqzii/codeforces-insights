import { SolvedTags, SolvedRatings, RatingHistory, hideLoader, showLoader, getRatingColor } from "./charts.js";
import { getPercentile, getPerformance, getRatingHistory, getSubmissions, getUserInfo } from "./api.js";
import { recommendProblems, setRatingRange, updateSubmissions } from "./recommendation.js";
import { setSearchValues } from "./search.js";
import { showError, toggleContentVisibility } from "./error.js"

const toggleOtherTags = document.getElementById("toggle-other-tags");
const toggle800Probs = document.getElementById("toggle-800-rating");
export const solvedTags = new SolvedTags(toggleOtherTags);
export const solvedRatings = new SolvedRatings(toggle800Probs);
export const ratingHistory = new RatingHistory();

const userDetailsEl = document.getElementById("user-details");
const solvedTagsEl = document.getElementById("solved-tags");
const solvedRatingsEl = document.getElementById("solved-ratings");
const ratingHistoryEl = document.getElementById("rating-history");

document.addEventListener("DOMContentLoaded", () => {
	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	toggleOtherTags.addEventListener("click", () => {
		solvedTags.toggleOther();
	});

	toggle800Probs.addEventListener("click", () => {
		solvedRatings.toggle800Rating();
	});
});

export async function updateAnalytics(handle, signal) {
	// Set charts to loading
	solvedTags.loading = true;
	toggleOtherTags.style.display = "none";
	solvedTags.updateChart();
	solvedRatings.loading = true;
	toggle800Probs.style.display = "none";
	solvedRatings.updateChart();

	ratingHistory.loading = true;
	ratingHistory.updateChart();
	showLoader(userDetailsEl);

	// Prevent displaying stale data.
	ratingHistory.updatePerfomanceData([]);
	ratingHistory.updateRatingData([]);
	ratingHistory.updateSolvedData([]);

	// Remove all previous error messages.
	const errorElems = document.querySelectorAll(".error-container");
	errorElems.forEach(el => {
		el.remove();
	});

	const userInfoTask = getUserInfo(handle, signal).then(info => {
		handleUserInfo(info);
		toggleContentVisibility(userDetailsEl, true);
		sessionStorage.setItem("userInfo", JSON.stringify(info));
	}).catch(err => {
		toggleContentVisibility(userDetailsEl, false);
		showError(err, userDetailsEl);
	});

	const submissionTask = getSubmissions(handle, signal).then(submissions => {
		const solved = filterSolved(submissions);
		handleSolved(solved);

		toggleContentVisibility(solvedTagsEl, true);
		toggleContentVisibility(solvedRatingsEl, true);

		updateSubmissions(submissions);
		sessionStorage.setItem("solvedProblems", JSON.stringify(solved));
	}).catch(err => {
		toggleContentVisibility(solvedTagsEl, false);
		showError(err, solvedTagsEl);
		toggleContentVisibility(solvedRatingsEl, false);
		showError(err, solvedRatingsEl);

		// Show error without hiding entire chart, as the rating history fetch can still succeed.
		showError(err, ratingHistoryEl);
	});

	getRatingHistory(handle, signal).then(ratings => {
		handleRatingHistory(handle, ratings, signal);
		toggleContentVisibility(ratingHistoryEl, true);
		sessionStorage.setItem("ratingHistory", JSON.stringify(ratings));
	}).catch(err => {
		toggleContentVisibility(ratingHistoryEl, false);
		showError(err, ratingHistoryEl);
	});

	// Recommend problems after we have user's rating and submissions.
	Promise.all([userInfoTask, submissionTask]).then(() => {
		recommendProblems();
	});
}

window.addEventListener("load", loadData);

function handleSolved(solved) {
	solved.sort((a, b) => {
		return a.creationTimeSeconds - b.creationTimeSeconds;
	});

	// Get count of each tag and rating
	const tagCnt = {}, ratingCnt = {};
	const solvedTime = new Array();
	solved.forEach(sub => {
		sub.problem.tags.forEach(tag => {
			tagCnt[tag] = (tagCnt[tag] || 0) + 1;
		});

		if (sub.problem.rating != undefined) {
			ratingCnt[sub.problem.rating] = (ratingCnt[sub.problem.rating] || 0) + 1;
			solvedTime.push({ timestamp: sub.creationTimeSeconds, rating: sub.problem.rating });
		}
	});

	const sortedTagCnt = Object.entries(tagCnt)
		.sort((a, b) => b[1] - a[1]);

	updateTags(sortedTagCnt);
	updateSolvedRatings(ratingCnt);
	updateSolvedRatingsTime(solvedTime);
}

function handleRatingHistory(handle, ratings, signal) {
	updateRatingChanges(ratings);

	const perfRequestData = new Array();
	ratings.forEach(rating => {
		perfRequestData.push({
			contestId: rating.contestId,
			oldRating: rating.oldRating,
			rank: rating.rank,
			ratingUpdateTimeSeconds: rating.ratingUpdateTimeSeconds,
		});
	});

	getPerformance(handle, perfRequestData, signal).then(perf => {
		updatePerformance(perf)
		sessionStorage.setItem("performance", JSON.stringify(perf));
	}).catch(err => {
		// Show error without hiding the other content, as we still have some information to display.
		showError(err, ratingHistoryEl);
	});
}

function handleUserInfo(userInfo, signal) {
	const rating = document.getElementById("user-rating");
	const peakRating = document.getElementById("user-peak-rating");
	const percentileElem = document.getElementById("user-percentile");
	if (userInfo.rating != undefined) {
		getPercentile(userInfo.rating, signal).then(percentile => {
			percentileElem.textContent = `${(percentile * 100).toFixed(2)}%`;
		}).catch(err => {
			showError(err, userDetailsEl);
		});
		percentileElem.classList.add("glow-text", "weight-600");

		rating.textContent = userInfo.rating;
		rating.style.setProperty("--text-color", getRatingColor(userInfo.rating));
		rating.classList.add("glow-color", "weight-450");

		peakRating.textContent = userInfo.maxRating;
		peakRating.style.setProperty("--text-color", getRatingColor(userInfo.maxRating));
		peakRating.classList.add("glow-color", "weight-450");

		// Update rating range for recommended problems
		const maxRatingDiff = 150;
		const maxRating = Math.max(1100, userInfo.rating + maxRatingDiff); // Recommend at least [800, 1100]
		setRatingRange(userInfo.rating - maxRatingDiff, maxRating);
	} else {
		percentileElem.textContent = "-";
		percentileElem.classList.remove("glow-text", "weight-600");

		rating.textContent = "-";
		rating.classList.remove("glow-color", "weight-450");

		peakRating.textContent = "-";
		peakRating.classList.remove("glow-color", "weight-450");

		setRatingRange(800, 1100); // Recommend low rated problems for user without rating.
	}

	hideLoader(userDetailsEl);
	document.getElementById("user-title-photo").src = userInfo.titlePhoto;
	document.getElementById("username").textContent = userInfo.handle;
	document.getElementById("username").href = "https://codeforces.com/profile/" + userInfo.handle;
	document.getElementById("user-country").textContent = userInfo.country || "-";

}

function filterSolved(submissions) {
	const solved = new Array();
	submissions.forEach(sub => {
		if (sub.verdict === "OK") {
			sub.problemId = sub.contestId + sub.problem.index;
			solved.push(sub);
		}
	});

	solved.sort((a, b) => a.problemId.localeCompare(b.problemId));
	const uniqueSolved = new Array();
	solved.forEach(sub => {
		if ((uniqueSolved.length == 0) || (sub.problemId != uniqueSolved.at(-1).problemId)) {
			uniqueSolved.push(sub)
		}
	});
	return uniqueSolved;
}

function updateTags(tagCnts) {
	solvedTags.updateData(tagCnts);
	solvedTags.loading = false;
	toggleOtherTags.style.display = "inline";
	solvedTags.updateChart();
}

function updateSolvedRatings(ratingCnts) {
	solvedRatings.updateData(ratingCnts);
	solvedRatings.loading = false;
	toggle800Probs.style.display = "inline";
	solvedRatings.updateChart();
}

function updateRatingChanges(ratingChanges) {
	ratingHistory.updateRatingData(ratingChanges);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

function updateSolvedRatingsTime(ratingsTime) {
	ratingHistory.updateSolvedData(ratingsTime);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

function updatePerformance(performance) {
	performance.sort((a, b) => a.timestamp - b.timestamp);

	ratingHistory.updatePerfomanceData(performance);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

function loadData() {
	const savedSolved = sessionStorage.getItem("solvedProblems");
	if (savedSolved) {
		const solved = JSON.parse(savedSolved);
		handleSolved(solved);
	}
	solvedTags.dataLoaded = true;
	solvedRatings.dataLoaded = true;

	const savedInfo = sessionStorage.getItem("userInfo");
	if (savedInfo) {
		const info = JSON.parse(savedInfo);
		handleUserInfo(info);

		setSearchValues(info.handle);
	}

	const savedRatings = sessionStorage.getItem("ratingHistory");
	if (savedRatings) {
		const ratings = JSON.parse(savedRatings);
		updateRatingChanges(ratings);
	}

	const savedPerformance = sessionStorage.getItem("performance");
	if (savedPerformance) {
		const perf = JSON.parse(savedPerformance);
		updatePerformance(perf);
	}

	ratingHistory.dataLoaded = true;
}
