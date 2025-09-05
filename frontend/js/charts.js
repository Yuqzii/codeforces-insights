import { fetchSolvedRatings } from "./api.js";

var solvedRatingsChart;

export async function updateSolvedRatingsChart(handle) {
	const ctx = document.getElementById('solved-ratings-chart');
	const data = await fetchSolvedRatings(handle);

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
			}
		}
	});
}
