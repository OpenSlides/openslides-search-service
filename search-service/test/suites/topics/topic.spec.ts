import { TestContainer } from '../../utils/test-container';
import { TestTopicRepositoryAdapter } from './test-topic-repository-adapter';

describe('Test searching for topics', () => {
    const testContainer = new TestContainer();
    const topicRepo = new TestTopicRepositoryAdapter(testContainer.datastore);

    beforeAll(() => {
        return testContainer.open();
    });

    afterAll(() => {
        return testContainer.close();
    });

    beforeEach(async () => {
        await testContainer.init();
    });

    it(`finds a topic`, async () => {
        const id = await topicRepo.create({ text: `Hello world!`, title: `Introduction` });
        const result = await testContainer.search.search(`hello`);
        console.log(`result:`, result.results);
    });
});
