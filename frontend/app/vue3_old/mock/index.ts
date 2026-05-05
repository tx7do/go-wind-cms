import categoryMock from './modules/category.mock'
import postMock from './modules/post.mock'
import commentMock from './modules/comment.mock'
import tagMock from './modules/tag.mock'
import pageMock from './modules/page.mock'
import navigationMock from './modules/navigation.mock'
import userMock from './modules/user.mock'

export default [
  ...categoryMock,
  ...postMock,
  ...commentMock,
  ...tagMock,
  ...pageMock,
  ...navigationMock,
  ...(Array.isArray(userMock) ? userMock : []),
];
