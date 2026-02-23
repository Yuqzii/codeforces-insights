# Endpoint: POST /performance
This endpoint is for calculating a users performance in a contest.


## Specification
### Request
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`handle`       |`String`  |Amount of problems to recommend. Must be in the range [0, 10].|
|`minRating`   |`Integer`  |Minimum possible rating of a user.|
|`maxRating`   |`Integer`  |Maximum possible rating of a user.|
|`ratingHistory`    |`Object[]` |Array of recent contests.|
|`contests[].id`|`Integer`  |The Codeforces contest ID.|
|`contests[].indices`|`String[]` |The indices of all problems solved for this contest.|

### Response
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`[].score`    |`Float`    |The similarity of the problem to the unsolved problems. Between -1 and 1.|
|`[].problem`  |`Problem`  |The recommended problem.|

## Flowchart
```mermaid
flowchart TD
	stadiumStart(["`Receive POST request to /performance`"])
	--> paraRevStart[\"Read request"\]
	--> WR["Calculate winrate for every possible rating difference using the standard elo rating for win probabability 1/(1+10^(delta_elo/400))"]
	--> Types["Calculate expected rank for every rating using the precalculated winrates, frequency of every rating amongst participants and Fast Fourier Transform"]
	--> Table["Make a table rankToRating out of the above process. Looking up a rating in this table yields its predicted rank"]
	--> DeclareLR["Begin binary search for rating-performance:<br>l := minRating + 1<br>r := maxRating"]
  --> BinarySearch{"l &lt; r?"}
  BinarySearch -- Yes --> DefMid["mid := (l + r)/2"]
  DefMid --> GetExpected["Look up expected rank for 'rating = mid' in the table rankToRating"]
  GetExpected --> Compare{"(Expected rank for mid) > (actual rank)?"}
  Compare -- Yes --> IncreaseL["l := mid + 1"]
  Compare -- No --> DecreaseR["r := mid"]
  IncreaseL --> BinarySearch
  DecreaseR --> BinarySearch
  --> paraRevEnd[\"`Write response: l - 1`"\]
  --> stadiumEnd(["`Send response`"])
```
