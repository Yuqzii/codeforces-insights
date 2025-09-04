const url = '/api/'

export async function fetchUserInfo(handle) {
	let endpoint = 'users/' + handle

	fetch(url + endpoint).then(response => {
		if (!response.ok)
			throw new Error('response not ok:', response.statusText)
		return response.json();
	}).then(data => {
		console.log(data)
	}).catch(error => {
		console.error(error)
	});
}
