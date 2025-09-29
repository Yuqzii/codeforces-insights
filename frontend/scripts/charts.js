var fgColor, bgColor, shadowColor, borderColor;
var redColor, orangeColor, greenColor, yellowColor, blueColor, purpleColor, aquaColor;

getColors();

const ratingRanges = [
	{ min: 0, max: 1199, color: '#ccccccff', label: "Newbie" },
	{ min: 1200, max: 1399, color: '#77ff77ff', label: "Pupil" },
	{ min: 1400, max: 1599, color: '#77ddbbff', label: "Specialist" },
	{ min: 1600, max: 1899, color: '#aaaaffff', label: "Expert" },
	{ min: 1900, max: 2099, color: '#ff88ffff', label: "Canditate Master" },
	{ min: 2100, max: 2299, color: '#ffcc88ff', label: "Master" },
	{ min: 2300, max: 2399, color: '#ffbb55ff', label: "International Master" },
	{ min: 2400, max: 2599, color: '#ff7777ff', label: "Grandmaster" },
	{ min: 2600, max: 2999, color: '#ff3333ff', label: "International Grandmaster" },
	{ min: 3000, max: 10000, color: '#aa0000ff', label: "Legendary Grandmaster" }
];

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

		const xAxisMin = this.#ratingData.labels[0];
		let xAxisMax = this.#ratingData.labels[this.#ratingData.labels.length - 1];
		if (this.#solvedData.length > 0)
			xAxisMax = Math.max(xAxisMax, this.#solvedData[this.#solvedData.length - 1].x);
		const xAxisPad = (xAxisMax - xAxisMin) * 0.01;

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
					borderColor: greenColor,
					backgroundColor: greenColor,
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
						min: xAxisMin - xAxisPad,
						max: xAxisMax + xAxisPad
					},
				},
				layout: {
					padding: {
						right: 10
					}
				},
				maintainAspectRatio: false,
				responsive: true
			},
			plugins: [{
				id: 'rankVisualPlugin',
				afterDraw: (chart) => {
					const { ctx, chartArea: { top, bottom, left, right }, scales: { y } } = chart;
					const lineWidth = 5;

					ctx.save();
					ctx.beginPath();
					ctx.rect(left - lineWidth, top, right - left + 2 * lineWidth, bottom - top);
					ctx.clip();

					ratingRanges.forEach(range => {
						const yMin = y.getPixelForValue(range.min);
						const yMax = y.getPixelForValue(range.max);

						ctx.save();
						ctx.beginPath();
						ctx.strokeStyle = range.color;
						ctx.lineWidth = lineWidth;
						ctx.moveTo(left - lineWidth / 2, yMin);
						ctx.lineTo(left - lineWidth / 2, yMax);
						ctx.moveTo(right + lineWidth / 2, yMin);
						ctx.lineTo(right + lineWidth / 2, yMax);
						ctx.stroke();


						ctx.restore();
					});

					ctx.restore();
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
