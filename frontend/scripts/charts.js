var solvedRatingsChart;

var fgColor, bgColor, shadowColor;
var redColor, orangeColor, greenColor, yellowColor, blueColor, purpleColor, aquaColor;
var grayDarkColor, blueDarkColor;

getColors();
Chart.defaults.color = fgColor;
Chart.defaults.borderColor = grayDarkColor;
Chart.defaults.datasets.bar.backgroundColor = blueColor;
Chart.defaults.elements.arc.backgroundColor = [redColor, greenColor, yellowColor, blueColor, purpleColor, orangeColor, aquaColor];

export class SolvedTags {
	N = 10;
	#showOtherTags = false;
	#tags = [];
	#counts = [];
	#chart;

	async updateChart() {
		const ctx = document.getElementById('solved-tags-chart');

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

		this.#chart = new Chart(ctx, {
			type: 'pie',
			data: {
				datasets: [{
					data: countsToShow
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
				responsive: true
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
	#chart;
	#data;

	updateChart() {
		const ctx = document.getElementById('solved-ratings-chart');

		if (this.#chart != null)
			this.#chart.destroy();

		hideLoader(ctx.parentNode.parentNode);
		solvedRatingsChart = new Chart(ctx, {
			type: 'bar',
			data: {
				datasets: [{
					label: '# of Solved Problems',
					data: this.#data,
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
			}
		});
	}

	updateData(data) {
		this.#data = data;
	}
}

function getColors() {
	var style = window.getComputedStyle(document.body);
	fgColor = style.getPropertyValue('--fg');
	bgColor = style.getPropertyValue('--bg');
	shadowColor = style.getPropertyValue('--shadow');
	redColor = style.getPropertyValue('--red');
	orangeColor = style.getPropertyValue('--orange');
	greenColor = style.getPropertyValue('--green');
	yellowColor = style.getPropertyValue('--yellow');
	blueColor = style.getPropertyValue('--blue');
	purpleColor = style.getPropertyValue('--purple');
	aquaColor = style.getPropertyValue('--aqua');
	grayDarkColor = style.getPropertyValue('--gray-dark');
	blueDarkColor = style.getPropertyValue('--blue-dark');
}

export function showLoader(container) {
	container.querySelector(".loader").style.display = "flex";
}

export function hideLoader(container) {
	container.querySelector(".loader").style.display = "none";
}
