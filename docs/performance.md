# Endpoint: POST /performance
This endpoint is for calculating a users performance in a contest.


## Specification
### Request
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`handle`       |`String`  |Codeforces handle of the user.|
|`minRating`   |`Integer`  |Minimum possible rating of a user.|
|`maxRating`   |`Integer`  |Maximum possible rating of a user.|
|`amountOfContests`|`Integer`  |Number of contests to be analyzed.|
|`ratingHistory`    |`Object[]` |Array of recent rating changes.|
|`contests`    |`Object[]` |Array of recent contests.|

### Response
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`performance`    |`integer[]`    |An array consisting of the given users performance in each given contest|


## Flowchart
```mermaid


flowchart TD


flowchart TD
	stadiumStart(["`Receive POST request to /performance`"])
	--> paraRevStart[\"Read request"\]
	--> WR["Calculate winrate for every possible rating difference using the standard elo rating for win probabability 1/(1+10^(delta_elo/400))"]
	--> iterator["i := 0"]
	--> forLoop["i < amountOfContests"]
	forLoop -- Yes	--> addJ*b["queue contest[i] as job"]
	--> increment[i := i+1] --> forLoop
	forLoop -- No --> waitForWorkers["Wait until all workers finish"]
	--> paraRevEnd[\"`Write response from result channel`"\]
	--> stadiumEnd(["`Send response`"])
	Worker(["initiate as worker"])
	--> WaitForJ*b["Wait until a job is available"] 
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
    BinarySearch -- No --> End["Send l-1 with contest timestamp to result channel"]
	--> WaitForJ*b
	
```
