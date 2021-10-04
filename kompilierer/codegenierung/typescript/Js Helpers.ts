export function eq(a: any, b: any): boolean {
    if (a === b) {
        return true
    }

    if (typeof a !== 'object' || a === null || b === null) {
        return false
    }

    for (var key in a) {
        if (!eq(a[key], b[key])) {
            return false
        }
    }

    return true
}