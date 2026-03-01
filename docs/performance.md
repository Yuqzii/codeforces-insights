# Endpoint: POST /performance
This endpoint is for calculating a users performance in any number of contests.


## Specification
### Request
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`handle`       |`String`  |Codeforces handle of the user.|
|`ratingHistory`    |`RatingChange` |Array of recent rating changes.|

### Response
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`[].Rating`    |`integer`    |The given users performance in the contest|
|`[].Timestamp`    |`integer`    |The contests timestamp|




## Flowchart
```mermaid
flowchart TD
	stadiumStart(["`Receive POST request to /performance`"])
	--> paraRevStart[\"Read request"\]
	--> WR["Calculate winrate for every possible rating difference using the standard elo rating for win probabability 1/(1+10^(delta_elo/400))"]
	--> iterator["i := 0 <br> jobQueue := queue <br> resultQueue := queue"]
	--> forLoop["i < amountOfContests"]
	forLoop -- Yes	--> addJ*b["queue contest[i] in jobQueue"]
	--> increment[i := i+1] --> forLoop
	forLoop -- No --> waitForWorkers["Wait until all workers finish"]
	--> paraRevEnd[\"`Read all of resultQueue and write  it as response`"\]
	--> stadiumEnd(["`Send response`"])

	Worker(["initiate as worker"])
	--> WaitForJ*b["Wait until a job is available in jobQueue"] 
	--> getJob["pop job from jobQueue"]
	--> Types["Calculate expected rank for every rating using the precalculated winrates, frequency of every rating amongst participants and Fast Fourier Transform"]
	--> Table["Make a table rankToRating out of the above process. Looking up a rating in this table yields its predicted rank"]
	--> DeclareLR["Begin binary search for rating performance:<br>l := minRating + 1<br>r := maxRating"]
    --> BinarySearch{"l &lt; r?"}
    BinarySearch -- Yes --> DefMid["mid := floor((l + r)/2)"]
    DefMid --> GetExpected["Look up expected rank for 'rating = mid' in the table rankToRating"]
    GetExpected --> Compare{"(Expected rank for mid) > (actual rank)?"}
    Compare -- Yes --> IncreaseL["l := mid + 1"]
    Compare -- No --> DecreaseR["r := mid"]
    IncreaseL --> BinarySearch
    DecreaseR --> BinarySearch
    BinarySearch -- No --> End["queue l-1 with contest timestamp to resultQueue]
	--> WaitForJ*b
```
