import { fetchSolvedTagsAndRatings } from "./api.js";

var solvedRatingsChart;
var solvedTagsChart;

export async function updateSolvedTagsAndRatingsCharts(handle) {
	const data = await fetchSolvedTagsAndRatings(handle);

	updateSolvedRatingsChart(data.ratings);
	updateSolvedTagsChart(data.tags);
}

async function updateSolvedRatingsChart(data) {
	const ctx = document.getElementById('solved-ratings-chart');

	if (solvedRatingsChart != null)
		solvedRatingsChart.destroy();

	solvedRatingsChart = new Chart(ctx, {
		type: 'bar',
		data: {
			datasets: [{
				label: '# of Solved Problems',
				data: data,
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

async function updateSolvedTagsChart(data) {
	const ctx = document.getElementById('solved-tags-chart');

	if (solvedTagsChart != null)
		solvedTagsChart.destroy();

	const tags = [];
	const counts = [];
	for (const element of data) {
		console.log(element);
		tags.push(element.tag);
		counts.push(element.count);
	}

	solvedTagsChart = new Chart(ctx, {
		type: 'pie',
		data: {
			datasets: [{
				data: counts
			}],
			labels: tags
		},
		options: {
			plugins: {
				legend: {
					display: false
				}
			},
			borderWidth: 1,
			responsive: true
		}
	});
}
