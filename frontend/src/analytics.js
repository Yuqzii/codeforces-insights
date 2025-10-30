import { SolvedTags, SolvedRatings, RatingHistory, hideLoader, showLoader, getRatingColor } from "./charts.js";
import { getPercentile, getPerformance, getRatingHistory, getSubmissions, getUserInfo } from "./api.js";

const toggleOtherTags = document.getElementById("toggle-other-tags");
const toggle800Probs = document.getElementById("toggle-800-rating");
export const solvedTags = new SolvedTags(toggleOtherTags);
export const solvedRatings = new SolvedRatings(toggle800Probs);
export const ratingHistory = new RatingHistory();

const userDetails = document.getElementById("user-details");

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
	showLoader(userDetails);

	// Prevent displaying stale data
	ratingHistory.updatePerfomanceData([]);
	ratingHistory.updateRatingData([]);
	ratingHistory.updateSolvedData([]);

	getUserInfo(handle, signal).then(handleUserInfo);
	getSubmissions(handle, signal).then(handleSubmissions);
	getRatingHistory(handle, signal).then(ratings => {
		handleRatingHistory(ratings, signal);
	});
}

function handleSubmissions(submissions) {
	submissions = filterSolved(submissions);
	submissions.sort((a, b) => {
		return a.creationTimeSeconds - b.creationTimeSeconds;
	});

	// Get count of each tag and rating
	const tagCnt = {}, ratingCnt = {};
	const solvedTime = new Array();
	submissions.forEach(sub => {
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

function handleRatingHistory(ratings, signal) {
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

	getPerformance(perfRequestData, signal).then(updatePerformance);
}

function handleUserInfo(userInfo, signal) {
	const rating = document.getElementById("user-rating");
	const peakRating = document.getElementById("user-peak-rating");
	const percentileElem = document.getElementById("user-percentile");
	if (userInfo.rating != undefined) {
		getPercentile(userInfo.rating, signal).then(percentile => {
			percentileElem.textContent = `${(percentile * 100).toFixed(2)}%`;
		});
		percentileElem.classList.add("glow-text", "weight-600");

		rating.textContent = userInfo.rating;
		rating.style.setProperty("--text-color", getRatingColor(userInfo.rating));
		rating.classList.add("glow-color", "weight-450");

		peakRating.textContent = userInfo.maxRating;
		peakRating.style.setProperty("--text-color", getRatingColor(userInfo.maxRating));
		peakRating.classList.add("glow-color", "weight-450");
	} else {
		percentileElem.textContent = "-";
		percentileElem.classList.remove("glow-text", "weight-600");

		rating.textContent = "-";
		rating.classList.remove("glow-color", "weight-450");

		peakRating.textContent = "-";
		peakRating.classList.remove("glow-color", "weight-450");
	}

	hideLoader(userDetails);
	document.getElementById("user-title-photo").src = userInfo.titlePhoto;
	document.getElementById("username").textContent = userInfo.handle;
	document.getElementById("user-country").textContent = userInfo.country || "-";

}

function filterSolved(submissions) {
	const solved = new Array();
	submissions.forEach(sub => {
		if (sub.verdict === "OK") solved.push(sub);
	});
	return solved;
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
	performance.sort((a, b) => a.timestamp > b.timestamp);

	ratingHistory.updatePerfomanceData(performance);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

