import { User } from '../../../src/domain/models/users/user';
import { TEST_ADMIN_USERNAME, TEST_ADMIN_PASSWORD, TEST_ADMIN_ID } from '../../utils/config';
import { Fqid } from '../../../src/domain/definitions/key-types';
import { RepositoryAdapter } from '../../utils/repository-adapter';

export class TestUserRepositoryAdapter extends RepositoryAdapter<User> {
    public createSuperadmin(): Promise<number> {
        return this.create({
            id: TEST_ADMIN_ID,
            username: TEST_ADMIN_USERNAME,
            default_password: TEST_ADMIN_PASSWORD,
            organization_management_level: `superadmin`
        });
    }

    protected getFqid(id: number): Fqid {
        return `${User.COLLECTION}/${id}`;
    }
}
