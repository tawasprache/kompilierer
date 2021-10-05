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
export function update(a: object, b: object): any {
    var c = {}

    for (var key in a) {
        c[key] = a[key]
    }

    for (var key in b) {
        c[key] = b[key]
    }

    return c
}