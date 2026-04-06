# Endpoint: POST /recommend
This endpoint is for recommending Codeforces problems similar to
problems the user has struggled with in earlier contests.


## Specification
### Request
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`count`       |`Integer`   |Amount of problems to recommend. Must be in the range [0, 10].|
|`minRating`   |`Integer`   |Minimum rating of recommended problems.|
|`maxRating`   |`Integer`   |Maximum rating of recommended problems.|
|`lookback`    |`Integer`   |How many recent contests to consider when making recommendations.|
|`submissions` |`Submission`|All accepted submissions the user has ever made.|

### Response
| Property     | Type      | Description |
|--------------|-----------|-------------|
|`[].score`    |`Float`    |The similarity of the problem to the unsolved problems. Between -1 and 1.|
|`[].problem`  |`Problem`  |The recommended problem.|

## Flowchart
```mermaid
flowchart TD
  A([Receive POST request to<br>/recommend])
  A --> B[/Read request as req/]
  B --> solvedSort
  subgraph Find solved problems from recent contests
    solvedSort["Sort submissions based on submission time"]
    solvedSort --> filterInit["probsByContest := map of contest ID to problem list<br>i := 0"]
    filterInit --> solvedLenCheck{"i < len(req.submissions)"}

    solvedLenCheck -- true --> isContestant{"Was submission i made as a contestant?"}
    isContestant -- false --> solvedInc
    isContestant -- true --> isStored{"probsByContest contains an entry for contest ID of submission i"}
    isStored -- true --> appendSubmission["Append problem of submission i to the list in the map"]
      appendSubmission --> solvedInc
    isStored -- false --> lookbackCheck{"len(probsByContest) < req.lookback"}
    lookbackCheck -- true --> newContestEntry["Create a new entry for the contest ID of submission i containing its corresponding problem"]
      newContestEntry --> solvedInc
    solvedInc["i += 1"] --> solvedLenCheck
  end

  subgraph Find first unsolved problem in recent contests
    solvedLenCheck -- false --> unsolvedInit
    unsolvedInit["unsolved := empty list of problems"]
    unsolvedInit --> unsolvedLenCheck{"probsByContest has remaining entries"}
    unsolvedLenCheck -- true --> getMapEntry["Get next (contest, problems) in probsByContest"]
      getMapEntry --> createIndexList["indices := sorted list of the indices of problems"]
      createIndexList --> getAllProbs[("allProblems := all problems of contest sorted by indices")]
      getAllProbs --> unsolvedFind["Find first problem in allProblems with index not in indices, and append to unsolved"]
      unsolvedFind --> unsolvedLenCheck
    
  end
  subgraph Find similar problems
    unsolvedLenCheck -- false --> I[Convert each problem in unsolved to vector]
    I --> J[target := sum of these vectors]
    J --> K[("Get all problems matching at least one tag of any problem in unsolved,
  and with rating in the range [req.minRating, req.MaxRating] as probs")]
    K --> L[i := 0<br>pq := prority queue of problems sorted on 'score']
    L --> M{"i < len(probs)"}
    M -- true --> N["convert probs[i] to vector v"]
    N --> O["similarity := dotProduct(target, v) / (target.magnitude() * v.magnitude())"]
    O --> P["Push probs[i] to pq with score as similarity"]
    P --> Q{"len(pq) > req.cnt"}
    Q -- true --> R["pq.pop()"] --> V
    Q -- false --> V[i++]
    V --> M
  end
  M -- false --> S[Marshal pq to json]
  S --> T[/Write json to response/]
  T --> U([Send response])
```
