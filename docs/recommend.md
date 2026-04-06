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
  B --> C
  subgraph Find unsolved problems
    C[i := 0<br>unsolved := empty array]
    C --> D{"i < len(req)"}
    D -- true --> E[("Get all problems of contests[i]")]
    E --> F["Sort problems based on index, e.g. ['A', 'B', 'C1']"]
    F --> G["Find first problem with index not in req.contests[i].indices, and append to unsolved"]
    G --> H[i += 1] --> D
  end
  subgraph Find similar problems
    D -- false --> I[Convert each problem in unsolved to vector]
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
