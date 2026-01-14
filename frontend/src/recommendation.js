export function filterSubmissionsRecentContests(submissions) {
    const filtered = [];
    const recentContests = [];
    
    for (let i = 0; i < submissions.length; i++){
        if (submissions[i].author.participantType ==  "CONTESTANT") {
            if (recentContests.includes(submissions[i].contestId)) {
                filtered.push(submissions[i]);
            } else if (recentContests.length < 5) {
                filtered.push(submissions[i]);
                recentContests.push(submissions[i].contestId);
            }
        }
    }
    
    return filtered
}
