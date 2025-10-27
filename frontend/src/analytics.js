import { SolvedTags, SolvedRatings, RatingHistory, hideLoader, showLoader, getRatingColor } from "./charts.js";
import { getPerformance, getRatingHistory, getSubmissions, getUserInfo } from "./codeforces.js";

const toggleOtherTags = document.getElementById('toggle-other-tags');
const toggle800Probs = document.getElementById('toggle-800-rating');
export const solvedTags = new SolvedTags(toggleOtherTags);
export const solvedRatings = new SolvedRatings(toggle800Probs);
export const ratingHistory = new RatingHistory();

const userDetails = document.getElementById('user-details');

document.addEventListener('DOMContentLoaded', () => {
	solvedTags.updateChart();
	solvedRatings.updateChart();
	ratingHistory.updateChart();

	toggleOtherTags.addEventListener('click', () => {
		solvedTags.toggleOther();
	});

	toggle800Probs.addEventListener('click', () => {
		solvedRatings.toggle800Rating();
	});
});


export async function updateAnalytics(handle, signal) {
	// Set charts to loading
	solvedTags.loading = true;
	toggleOtherTags.style.display = 'none';
	solvedTags.updateChart();
	solvedRatings.loading = true;
	toggle800Probs.style.display = 'none';
	solvedRatings.updateChart();

	ratingHistory.loading = true;
	ratingHistory.updateChart();
	showLoader(userDetails);

	// Prevent displaying stale data
	ratingHistory.updatePerfomanceData([]);
	ratingHistory.updateRatingData([]);
	ratingHistory.updateSolvedData([]);

	getUserInfo(handle, (userInfo) => {
		updateUserInfo(userInfo);
	}, signal);

	// Update everything that needs submission history
	getSubmissions(handle, (submissions) => {
		submissions = filterSolved(submissions);

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
	}, signal);

	getRatingHistory(handle, (ratings) => {
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

		getPerformance(perfRequestData, (performance) => {
			console.log(performance);
			updatePerformance(performance);
		}, signal);
	});
}

function filterSolved(submissions) {
	const solved = new Array();
	submissions.forEach(sub => {
		if (sub.verdict === "OK") solved.push(sub);
	});
	return solved;
}

async function updateTags(tagCnts) {
	solvedTags.updateData(tagCnts);
	solvedTags.loading = false;
	toggleOtherTags.style.display = 'inline';
	solvedTags.updateChart();
}

async function updateSolvedRatings(ratingCnts) {
	solvedRatings.updateData(ratingCnts);
	solvedRatings.loading = false;
	toggle800Probs.style.display = 'inline';
	solvedRatings.updateChart();
}

async function updateRatingChanges(ratingChanges) {
	ratingHistory.updateRatingData(ratingChanges);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

async function updateSolvedRatingsTime(ratingsTime) {
	ratingHistory.updateSolvedData(ratingsTime);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

async function updatePerformance(performance) {
	performance.sort((a, b) => a.timestamp > b.timestamp);

	ratingHistory.updatePerfomanceData(performance);
	ratingHistory.loading = false;
	ratingHistory.updateChart();
}

async function updateUserInfo(userInfo) {
	hideLoader(userDetails);
	document.getElementById('user-title-photo').src = userInfo.titlePhoto;
	document.getElementById('username').textContent = userInfo.handle;
	document.getElementById('user-country').textContent = userInfo.country;

	const rating = document.getElementById('user-rating');
	rating.textContent = userInfo.rating;
	rating.style.setProperty('--text-color', getRatingColor(userInfo.rating));
	const peakRating = document.getElementById('user-peak-rating');
	peakRating.textContent = userInfo.maxRating;
	peakRating.style.setProperty('--text-color', getRatingColor(userInfo.maxRating));
}
