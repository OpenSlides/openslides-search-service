import { Client } from 'pg';
import { PostgreAdapterService, PostgreService } from '../../src/gateways/postgre';

const ALL_TABLES = [
    'positions',
    'events',
    'id_sequences',
    'collectionfields',
    'events_to_collectionfields',
    'models',
    'migration_keyframes',
    'migration_keyframe_models',
    'migration_events',
    'migration_positions'
];

const ALL_SEQUENCES = ['positions_position', 'events_id', 'collectionfields_id'];

export class TestPostgreAdapter {
    private readonly _client: PostgreService = new PostgreAdapterService();

    public async prune(): Promise<void> {
        try {
            const client = await this.getClient();
            for (const table of ALL_TABLES) {
                await client.query(`DELETE FROM ${table} CASCADE;`, []);
            }
            for (const sequence of ALL_SEQUENCES) {
                await client.query(`ALTER SEQUENCE ${sequence}_seq RESTART WITH 1;`, []);
            }
        } catch (e: any) {
            console.log('Error prune', e.stack);
        }
    }

    public async openConnection(): Promise<void> {
        await this.getClient();
    }

    public async closeConnection(): Promise<void> {
        await (await this.getClient()).end();
    }

    private getClient(): Promise<Client> {
        return this._client.getPgClient() as Promise<Client>;
    }
}
