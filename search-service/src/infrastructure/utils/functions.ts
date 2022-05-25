import { RestServerResponse } from './definitions';

interface CreateResponseConfig<T = void> {
    results?: T[];
    message?: string;
    success?: boolean;
}

export const createResponse = <T = unknown>({
    results = [],
    message = 'Action handled successfully',
    success = true
}: CreateResponseConfig<T> = {}): RestServerResponse<T> => ({ message, success, results });
