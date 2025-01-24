export function k_combinations(set, k) {
    const combs = []
    let j, head, tailcombs;

    if (k > set.length || k <= 0) {
        return [];
    }

    if (k == set.length) {
        return [set];
    }

    if (k == 1) {
        for (let i = 0; i < set.length; i++) {
            combs.push([set[i]]);
        }
        return combs;
    }

    for (let i = 0; i < set.length - k + 1; i++) {
        head = set.slice(i, i + 1);
        tailcombs = k_combinations(set.slice(i + 1), k - 1);
        for (j = 0; j < tailcombs.length; j++) {
            combs.push(head.concat(tailcombs[j]));
        }
    }

    return combs;
}