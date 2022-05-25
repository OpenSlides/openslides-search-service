export const DATASTORE_DATABASE_HOST = process.env.DATASTORE_DATABASE_HOST as string;
export const SEARCH_SERVICE_PORT = parseInt(process.env.SEARCH_PORT || '', 10) || 9022;
export const SEARCH_SERVICE_HOST = process.env.INSTANCE_DOMAIN || 'http://localhost';

export class Config {
    public static readonly DATABASE_HOST = DATASTORE_DATABASE_HOST;
    public static readonly PORT: number = SEARCH_SERVICE_PORT;
    public static readonly DOMAIN: string = SEARCH_SERVICE_HOST;

    private static readonly VERBOSE_TRUE_FIELDS = ['1', 'y', 'yes', 'true', 'on'];

    public static isDevMode(): boolean {
        return this.VERBOSE_TRUE_FIELDS.includes(process.env.OPENSLIDES_DEVELOPMENT || '');
    }
}
