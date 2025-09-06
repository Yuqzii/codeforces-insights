import { fetchSolvedTagsAndRatings } from "./api.js";
import { updateSolvedTagsChart, updateTagsChartData } from "./solvedTags.js";

var solvedRatingsChart;

export async function updateSolvedTagsAndRatingsCharts(handle) {
	const data = await fetchSolvedTagsAndRatings(handle);

	updateSolvedRatingsChart(data.ratings);

	updateTagsChartData(data.tags);
	updateSolvedTagsChart();
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

