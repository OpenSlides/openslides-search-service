declare global {
    /**
     * Enhance array with own functions
     */
    interface Array<T> {
        mapToObject<U>(f: (item: T, index: number) => { [key in keyof U]: any }): { [key in keyof U]: any };
        /**
         * TODO: Remove this, when ES 2019 (or in whatever spec `at` is defined) is the target for our tsconfig
         *
         * @param index The index of an element in the array one expects
         */
        at(index: number): T | undefined;
    }

    /**
     * Enhances the number object to calculate real modulo operations.
     * (not remainder)
     */
    interface Number {
        modulo(n: number): number;
    }
}

export function overloadJsFunctions(): void {
    overloadArrayFunctions();
    overloadNumberFunctions();
}

function overloadArrayFunctions(): void {
    Object.defineProperty(Array.prototype, `toString`, {
        value(): string {
            let string = ``;
            const iterations = Math.min(this.length, 3);

            for (let i = 0; i <= iterations; i++) {
                if (i < iterations) {
                    string += this[i];
                }

                if (i < iterations - 1) {
                    string += `, `;
                } else if (i === iterations && this.length > iterations) {
                    string += `, ...`;
                }
            }
            return string;
        },
        enumerable: false
    });

    Object.defineProperty(Array.prototype, `mapToObject`, {
        value<T, U extends object>(f: (item: T, index: number) => U): U {
            return this.reduce((aggr: U, item: T, index: number) => {
                const res = f(item, index);
                for (const key in res) {
                    if (res.hasOwnProperty(key)) {
                        aggr[key] = res[key];
                    }
                }
                return aggr;
            }, {});
        },
        enumerable: false
    });

    Object.defineProperty(Array.prototype, `at`, {
        value<T>(index: number): T | undefined {
            if (index < 0) {
                index = index.modulo(this.length);
            }
            if (index > this.length) {
                return undefined;
            }
            return this[index];
        },
        enumerable: false
    });
}

/**
 * Enhances the number object with a real modulo operation (not remainder).
 * TODO: Remove this, if the remainder operation is changed to modulo.
 */
function overloadNumberFunctions(): void {
    Object.defineProperty(Number.prototype, `modulo`, {
        value(n: number): number {
            return ((this % n) + n) % n;
        },
        enumerable: false
    });
}
