export function filterSubmissionsRecentContests(submissions,amountOfContests) {
    const filtered = [];
    const recentContests = [];

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

    return  {submissions: filtered, contestId: recentContests}
}

export function findSolvedProblemsRecentContests(contestSubmissions) {
    const solvedProblems = [];

    for (const contest of contestSubmissions.submissions) {
        const solvedProblemsInAContest = [];
        for (const problem of contest)
            if (problem.verdict=="OK") {
                solvedProblemsInAContest.push(problem.problem.index);
            }
        solvedProblems.push({id:contestSubmissions.contestId.shift(), indices:solvedProblemsInAContest})
    }

    return solvedProblems
}
