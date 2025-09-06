const N = 10;

let showOtherTags = false;
var chart;
let tags = [];
let counts = [];

export async function updateSolvedTagsChart() {
	const ctx = document.getElementById('solved-tags-chart');

	let tagsToShow = [];
	let countsToShow = [];
	if (showOtherTags) {
		tagsToShow = tags;
		countsToShow = counts;
	} else {
		// Display top N tags
		tagsToShow = tags.slice(0, N);
		countsToShow = counts.slice(0, N);

		tagsToShow.push("Other");
		let otherCount = 0;
		for (let i = N; i < counts.length; i++)
			otherCount += counts[i]
		countsToShow.push(otherCount);
	}

	if (chart != null)
		chart.destroy();

	chart = new Chart(ctx, {
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

export function updateTagsChartData(data) {
	data.reverse();
	tags = [];
	counts = [];
	for (const element of data) {
		tags.push(element.tag);
		counts.push(element.count);
	}
}

export function toggleShowOtherTags() {
	showOtherTags = !showOtherTags;
	updateSolvedTagsChart();
}
