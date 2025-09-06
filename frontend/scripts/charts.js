import { fetchSolvedTagsAndRatings } from "./api.js";
import { updateSolvedTagsChart, updateTagsChartData } from "./solvedTags.js";

var solvedRatingsChart;

var fgColor, bgColor, shadowColor;
var redColor, orangeColor, greenColor, yellowColor, blueColor, purpleColor, aquaColor;
var grayDarkColor, blueDarkColor;

await getColors();
Chart.defaults.borderColor = grayDarkColor;
Chart.defaults.datasets.bar.backgroundColor = blueColor;
Chart.defaults.elements.arc.backgroundColor = [redColor, greenColor, yellowColor, blueColor, purpleColor, orangeColor, aquaColor];

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

async function getColors() {
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
