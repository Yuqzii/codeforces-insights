export function filterSubmissionsRecentContests(submissions,amountOfContests) {
    const filtered = [];
    const recentContests = [];
    const returnObject = [];

    for (let i = 0; i < submissions.length; i++){
        if (submissions[i].author.participantType != "CONTESTANT") {
            continue;
        }

        let ID=submissions[i].contestId;
        if (recentContests.includes(ID)) {
            filtered[recentContests.indexOf(ID)].push(submissions[i]);
        } else if (recentContests.length < amountOfContests) {
            filtered.push([submissions[i]]);
            recentContests.push(ID);
        }
    }

    for (let i = 0; i < recentContests.length; i++) {
        returnObject.push({id: recentContests[i], submissions: filtered[i]})
    }

    return  returnObject
}

export function findSolvedProblemsRecentContests(contests) {
    const solvedByContests = [];

    for (const contest of contests) {
        const solvedInContest = [];

        for (const problem of contest.submissions) {
            if (problem.verdict=="OK" && !(solvedInContest.includes(problem.index))) {
                solvedInContest.push(problem.problem.index);
            }
        }

        solvedByContests.push({id:contest.id, indices:solvedInContest})
    }

    return solvedByContests
}
