const routes = [
  {
    path: '/login',
    component: './Login',
    layout: false,
  },
  {
    path: '/',
    redirect: '/home',
  },
  {
    name: '代理',
    path: '/home',
    component: './Home',
  },
  {
    name: '应用/设备',
    path: '/device',
    routes: [
      {
        path: '/device',
        redirect: '/device/app',
      },
      {
        name: '应用',
        path: '/device/app',
        component: './AppList',
      },
      {
        name: '设备',
        path: '/device/list',
        component: './Device',
      },
    ],
  },
  {
    name: '连接器',
    path: '/access',
    component: './Access',
  },
  {
    name: '设置',
    path: '/table',
    component: './Table',
  },
  {
    name: '个人中心',
    path: '/profile',
    component: './Profile',
    hideInMenu: true,
  },
  {
    name: '账户管理',
    path: '/users',
    component: './UserManagement',
  },
];

export default routes;
