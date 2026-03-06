const url = "https://codeforces.com/api/"

async function cfFetch(endpoint, signal) {
	try {
		const resp = await fetch(url + endpoint, { signal });
		const data = await resp.json();
		if (data.status !== "OK") throw new Error(`Codeforces not OK: ${data.comment}`);
		return data.result;
	} catch (err) {
		if (err.name === "AbortError") return;
		console.error("Codeforces request failed:", err);
		throw err;
	}
}

export async function getUserInfo(handle, signal) {
	const data = await cfFetch(`user.info?handles=${handle}`, signal);
	return data[0];
}

export async function getSubmissions(handle, signal) {
	return await cfFetch(`user.status?handle=${handle}`, signal);
}

export async function getRatingHistory(handle, signal) {
	return await cfFetch(`user.rating?handle=${handle}`, signal);
}

export async function getProblems(tags, signal) {
	let url = "problemset.problems";
	if (tags) {
		url += "?tags=";
		for (const tag of tags) {
			url += tag + ";";
		}
		url = url.slice(0, -1);
		console.log(url)
	}

	return await cfFetch(url, signal);
}

export async function getPerformance(handle, ratingHistory, signal) {
	try {
		const reqData = {
			handle: handle,
			ratingHistory: ratingHistory
		};
		const resp = await fetch(`${process.env.API_URL}/performance`, {
			method: "POST",
			body: JSON.stringify(reqData),
			signal: signal,
		});
		if (!resp.ok) throw new Error(`performance response not ok: ${await resp.text()}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		throw err;
	}
}

export async function getPercentile(rating, signal) {
	try {
		const resp = await fetch(`${process.env.API_URL}/percentile/${rating}`, { signal });
		if (!resp.ok) throw new Error(`percentile response not ok: ${await resp.text()}`);
		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		throw err;
	}
}

export async function getRecommendedProblems(count, ratingRange, lookback, solved, signal) {
	const subs = minifySubmissions(solved);

	const reqData = {
		count: count,
		minRating: ratingRange.min,
		maxRating: ratingRange.max,
		lookback: lookback,
		submissions: subs,
	};

	try {
		const resp = await fetch(`${process.env.API_URL}/recommend`, {
			method: "POST",
			body: JSON.stringify(reqData),
			signal: signal,
		});

		if (!resp.ok)
			throw new Error(`recommend response not ok: ${await resp.text()}`);

		const data = await resp.json();
		return data;
	} catch (err) {
		if (err.name === "AbortError") return;
		throw err;
	}
}

// Removes unnecessary properties of submissions for more efficient data transfer.
function minifySubmissions(subs) {
	const res = new Array();

	subs.forEach(sub => {
		// Only exactly what we store on the backend.
		res.push({
			id: sub.id,
			contestId: sub.contestId,
			verdict: sub.verdict,
			problem: {
				name: sub.problem.name,
				contestId: sub.problem.contestId,
				index: sub.problem.index,
				rating: sub.problem.rating,
				tags: sub.problem.tags,
			},
			author: {
				participantType: sub.author.participantType,
			},
			programmingLanguage: sub.programmingLanguage,
			creationTimeSeconds: sub.creationTimeSeconds,
		});
	});

	return res;
}
