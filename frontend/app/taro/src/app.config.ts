export default defineAppConfig({
  pages: [
    'pages/index',
    'pages/about/index',
    'pages/contact/index',
    'pages/login/index',
    'pages/register/index',
    'pages/post/[id]/index',
    'pages/category/[id]/index',
    'pages/tag/index',
    'pages/tag/[id]/index',
    'pages/settings/index',
    'pages/user/index'
  ],
  window: {
    backgroundTextStyle: 'light',
    navigationBarBackgroundColor: '#fff',
    navigationBarTitleText: 'Taro CMS',
    navigationBarTextStyle: 'black'
  }
})
