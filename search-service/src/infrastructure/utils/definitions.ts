export interface RestServerResponse<T = unknown> {
    results: T[];
    message: string;
    success: boolean;
}
