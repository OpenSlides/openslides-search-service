export interface Constructable<T = any> {
    new (...args: any[]): T;
    prototype: T;
    name?: string;
}
