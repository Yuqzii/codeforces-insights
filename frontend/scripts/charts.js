var fgColor, bgColor, shadowColor, borderColor;
var redColor, orangeColor, greenColor, yellowColor, blueColor, purpleColor, aquaColor;

const newbieColor = '#cccccc55', pupilColor = '#77ff7775', specialistColor = '#77ddbb75', expertColor = '#aaaaff75',
	cmasterColor = '#ff88ff75', masterColor = '#ffcc8875', imasterColor = '#ffbb5575', gmasterColor = '#ff777775',
	igmasterColor = '#ff333375', lgmasterColor = '#aa000075';

getColors();

export class SolvedTags {
	N = 10;
	loading = true;
	#showOtherTags = false;
	#tags = [];
	#counts = [];
	#chart;

	constructor(toggleOtherButton) {
		this.toggleOtherButton = toggleOtherButton
	}

	async updateChart() {
		const ctx = document.getElementById('solved-tags-chart');

		if (this.loading) {
			showLoader(ctx.parentNode.parentNode);
			return;
		}

		let tagsToShow = [];
		let countsToShow = [];
		if (this.#showOtherTags) {
			tagsToShow = this.#tags;
			countsToShow = this.#counts;
		} else {
			// Display top N tags
			tagsToShow = this.#tags.slice(0, this.N);
			countsToShow = this.#counts.slice(0, this.N);

			tagsToShow.push("Other");
			let otherCount = 0;
			for (let i = this.N; i < this.#counts.length; i++)
				otherCount += this.#counts[i]
			countsToShow.push(otherCount);
		}

		if (this.#chart != null)
			this.#chart.destroy();

		hideLoader(ctx.parentNode.parentNode);

		this.toggleOtherButton.style.display = 'inline';

		this.#chart = new Chart(ctx, {
			type: 'pie',
			data: {
				datasets: [{
					data: countsToShow,
					backgroundColor: [redColor, greenColor, yellowColor, blueColor, purpleColor, orangeColor, aquaColor]
				}],
				labels: tagsToShow
			},
			options: {
				plugins: {
					legend: {
						display: false
					}
				},
				borderWidth: 0.5,
				responsive: true,
				maintainAspectRatio: true
			}
		});
	}

	updateData(data) {
		data.reverse();
		this.#tags = [];
		this.#counts = [];
		for (const element of data) {
			this.#tags.push(element.tag);
			this.#counts.push(element.count);
		}
	}

	toggleOther() {
		this.#showOtherTags = !this.#showOtherTags;
		this.updateChart();
	}
}

export class SolvedRatings {
	loading = true;
	#chart;
	#data;

	updateChart() {
		const ctx = document.getElementById('solved-ratings-chart');

		if (this.loading) {
			showLoader(ctx.parentNode.parentNode);
			return;
		}

		if (this.#chart != null)
			this.#chart.destroy();

		hideLoader(ctx.parentNode.parentNode);
		this.#chart = new Chart(ctx, {
			type: 'bar',
			data: {
				datasets: [{
					label: '# of Solved Problems',
					data: this.#data,
					backgroundColor: blueColor,
				}]
			},
			options: {
				scales: {
					y: {
						beginAtZero: true
					}
				},
				elements: {
					bar: {
						borderRadius: 8
					}
				},
				maintainAspectRatio: false,
				responsive: true
			},
		});
	}

	updateData(data) {
		this.#data = data;
	}
}

export class RatingHistory {
	loading = true;
	#chart;
	#ratingData = [new Array, new Array];
	#performanceData = [new Array, new Array];
	#solvedData = new Array;

	updateChart() {
		const ctx = document.getElementById('rating-history-chart');

		if (this.loading) {
			showLoader(ctx.parentNode.parentNode);
			return;
		}

		if (this.#chart != null)
			this.#chart.destroy();

		hideLoader(ctx.parentNode.parentNode);
		this.#chart = new Chart(ctx, {
			type: 'line',
			data: {
				labels: this.#ratingData.labels,
				datasets: [{
					label: 'Rating',
					data: this.#ratingData.ratings,
					tension: 0.25,
					borderColor: orangeColor,
					backgroundColor: orangeColor,
				}, {
					label: 'Performance',
					data: this.#performanceData.performance,
					tension: 0.25,
					borderColor: aquaColor,
					backgroundColor: aquaColor,
				},
				{
					label: 'Solved Problems',
					type: 'scatter',
					data: this.#solvedData,
					borderColor: blueColor,
					backgroundColor: blueColor,
					pointRadius: 2.5,
				}]
			},
			options: {
				responsive: true,
				scales: {
					x: {
						type: 'time',
						time: {
							unit: 'month'
						},
						min: this.#ratingData.labels[0],
						max: Math.max(this.#ratingData.labels[this.#ratingData.labels.length - 1],
							this.#solvedData[this.#solvedData.length - 1]),
					}
				},
				maintainAspectRatio: false,
				responsive: true
			},
			plugins: [{
				id: 'backgroundColorByY',
				beforeDraw: (chart) => {
					const { ctx, chartArea: { top, bottom, left, right }, scales: { y } } = chart;

					// Define y-ranges with colors
					const ranges = [
						{ from: 0, to: 1200, color: newbieColor },
						{ from: 1200, to: 1400, color: pupilColor },
						{ from: 1400, to: 1600, color: specialistColor },
						{ from: 1600, to: 1900, color: expertColor },
						{ from: 1900, to: 2100, color: cmasterColor },
						{ from: 2100, to: 2300, color: masterColor },
						{ from: 2300, to: 2400, color: imasterColor },
						{ from: 2400, to: 2600, color: gmasterColor },
						{ from: 2600, to: 3000, color: igmasterColor },
						{ from: 3000, to: 10000, color: lgmasterColor }
					];

					ranges.forEach(range => {
						const yStart = y.getPixelForValue(range.to);
						const yEnd = y.getPixelForValue(range.from);

						ctx.fillStyle = range.color;
						ctx.fillRect(left, yStart, right - left, yEnd - yStart);
					});
				}
			}]
		});
	}

	updateRatingData(data) {
		this.#ratingData.ratings = [];
		this.#ratingData.labels = [];
		for (const element of data) {
			this.#ratingData.ratings.push(element.newRating);
			this.#ratingData.labels.push(element.ratingUpdateTimeSeconds * 1000);
		}
	}

	updatePerfomanceData(data) {
		this.#performanceData.performance = [];
		this.#performanceData.timestamps = [];
		for (const element of data) {
			this.#performanceData.performance.push(element.rating);
			this.#performanceData.timestamps.push(element.timestamp);
		}
	}

	updateSolvedData(data) {
		this.#solvedData = data.map(el => ({
			x: el.timestamp * 1000,
			y: el.rating
		}));
	}
}

export function getColors() {
	const style = window.getComputedStyle(document.documentElement);
	fgColor = style.getPropertyValue('--fg');
	bgColor = style.getPropertyValue('--bg');
	shadowColor = style.getPropertyValue('--shadow');
	borderColor = style.getPropertyValue('--border');
	redColor = style.getPropertyValue('--red');
	orangeColor = style.getPropertyValue('--orange');
	greenColor = style.getPropertyValue('--green');
	yellowColor = style.getPropertyValue('--yellow');
	blueColor = style.getPropertyValue('--blue');
	purpleColor = style.getPropertyValue('--purple');
	aquaColor = style.getPropertyValue('--aqua');

	Chart.defaults.color = fgColor;
	Chart.defaults.borderColor = borderColor;
}

export function showLoader(container) {
	container.querySelector(".loader").style.display = "flex";
}

export function hideLoader(container) {
	container.querySelector(".loader").style.display = "none";
}
