/* eslint-disable @typescript-eslint/no-unused-vars */
const users = [
  { id: 0, name: 'Umi', nickName: 'U', gender: 'MALE' },
  { id: 1, name: 'Fish', nickName: 'B', gender: 'FEMALE' },
];

// 用户管理列表数据
let userList = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    avatar:
      'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
    role: 'admin',
    created_at: '2024-01-01 00:00:00',
    updated_at: '2024-01-01 00:00:00',
  },
  {
    id: 2,
    username: 'user1',
    email: 'user1@example.com',
    avatar: '',
    role: 'user',
    created_at: '2024-01-02 10:00:00',
    updated_at: '2024-01-02 10:00:00',
  },
  {
    id: 3,
    username: 'user2',
    email: 'user2@example.com',
    avatar: '',
    role: 'user',
    created_at: '2024-01-03 15:30:00',
    updated_at: '2024-01-03 15:30:00',
  },
];

// 应用列表数据 - 匹配swagger接口结构
const appList = [
  {
    id: 0,
    name: 'TradeCode 0',
    application_type: 'web',
    ip: '192.168.1.100',
    port: 357,
    edge_id: 1,
    device: {
      id: 1,
      name: '设备1',
      created_at: '2017-10-31 23:12:00',
      updated_at: '2017-10-31 23:12:00',
    },
    created_at: '2017-10-31 23:12:00',
    updated_at: '2017-10-31 23:12:00',
  },
  {
    id: 1,
    name: 'TradeCode 1',
    application_type: 'api',
    ip: '192.168.1.101',
    port: 944,
    edge_id: 1,
    device: {
      id: 2,
      name: '设备2',
      created_at: '2017-10-31 24:15:00',
      updated_at: '2017-10-31 24:15:00',
    },
    created_at: '2017-10-31 24:15:00',
    updated_at: '2017-10-31 24:15:00',
  },
  {
    id: 2,
    name: 'TradeCode 2',
    application_type: 'service',
    ip: '192.168.1.102',
    port: 594,
    edge_id: 2,
    device: {
      id: 3,
      name: '设备3',
      created_at: '2017-10-31 24:27:02',
      updated_at: '2017-10-31 24:27:02',
    },
    created_at: '2017-10-31 24:27:02',
    updated_at: '2017-10-31 24:27:02',
  },
  {
    id: 3,
    name: 'TradeCode 3',
    application_type: 'web',
    ip: '192.168.1.103',
    port: 239,
    edge_id: 2,
    device: {
      id: 4,
      name: '设备4',
      created_at: '2017-10-31 00:00:02',
      updated_at: '2017-10-31 00:00:02',
    },
    created_at: '2017-10-31 00:00:02',
    updated_at: '2017-10-31 00:00:02',
  },
  {
    id: 4,
    name: 'TradeCode 4',
    application_type: 'api',
    ip: '192.168.1.104',
    port: 988,
    edge_id: 3,
    device: {
      id: 5,
      name: '设备5',
      created_at: '2017-10-31 00:00:02',
      updated_at: '2017-10-31 00:00:02',
    },
    created_at: '2017-10-31 00:00:02',
    updated_at: '2017-10-31 00:00:02',
  },
];

export default {
  'GET /api/v1/queryUserList': (req: any, res: any) => {
    res.json({
      success: true,
      data: { list: users },
      errorCode: 0,
    });
  },
  'PUT /api/v1/user/': (req: any, res: any) => {
    res.json({
      success: true,
      errorCode: 0,
    });
  },
  // 应用列表接口 - 匹配swagger接口
  'GET /v1/applications': (req: any, res: any) => {
    const { page = 1, page_size = 10, name, status } = req.query;

    let filteredList = [...appList];

    // 根据应用名过滤
    if (name) {
      filteredList = filteredList.filter((item) =>
        item.name.toLowerCase().includes(name.toLowerCase()),
      );
    }

    // 根据状态过滤
    if (status) {
      filteredList = filteredList.filter((item) => item.status === status);
    }

    // 分页
    const start = (page - 1) * page_size;
    const end = start + page_size;
    const applications = filteredList.slice(start, end);

    res.json({
      code: 0,
      message: 'success',
      data: {
        applications,
        total: filteredList.length,
      },
    });
  },
  // IAM 登录接口
  'POST /v1/iam/login': (req: any, res: any) => {
    const { username, password } = req.body;
    if (username === 'admin' && password === 'admin') {
      res.json({
        code: 200,
        message: 'success',
        data: {
          token: 'mock-jwt-token-' + Date.now(),
          user: {
            id: 1,
            username: 'admin',
            email: 'admin@example.com',
            avatar:
              'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
            role: 'admin',
            created_at: '2024-01-01 00:00:00',
            updated_at: '2024-01-01 00:00:00',
          },
        },
      });
    } else {
      res.status(401).json({
        code: 401,
        message: '用户名或密码错误',
      });
    }
  },
  // IAM 登出接口
  'POST /v1/iam/logout': (req: any, res: any) => {
    res.json({
      code: 200,
      message: '登出成功',
    });
  },
  // IAM 获取用户信息接口
  'GET /v1/iam/profile': (req: any, res: any) => {
    const token = req.headers.authorization;
    if (token && token.startsWith('Bearer ')) {
      res.json({
        code: 200,
        message: 'success',
        data: {
          id: 1,
          username: 'admin',
          email: 'admin@example.com',
          avatar:
            'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
          role: 'admin',
          created_at: '2024-01-01 00:00:00',
          updated_at: '2024-01-01 00:00:00',
        },
      });
    } else {
      res.status(401).json({
        code: 401,
        message: '未登录或登录已过期',
      });
    }
  },
  // IAM 更新用户信息接口
  'PUT /v1/iam/profile': (req: any, res: any) => {
    const token = req.headers.authorization;
    if (token && token.startsWith('Bearer ')) {
      const { username, email, avatar } = req.body;
      res.json({
        code: 200,
        message: 'success',
        data: {
          id: 1,
          username: username || 'admin',
          email: email || 'admin@example.com',
          avatar:
            avatar ||
            'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
          role: 'admin',
          created_at: '2024-01-01 00:00:00',
          updated_at: new Date().toISOString(),
        },
      });
    } else {
      res.status(401).json({
        code: 401,
        message: '未登录或登录已过期',
      });
    }
  },
  // IAM 获取用户列表接口（管理员）
  'GET /v1/iam/users': (req: any, res: any) => {
    const { page = 1, page_size = 10, username, email, role } = req.query;

    let filteredList = [...userList];

    // 根据用户名过滤
    if (username) {
      filteredList = filteredList.filter((item) =>
        item.username.toLowerCase().includes(username.toLowerCase()),
      );
    }

    // 根据邮箱过滤
    if (email) {
      filteredList = filteredList.filter((item) =>
        item.email.toLowerCase().includes(email.toLowerCase()),
      );
    }

    // 根据角色过滤
    if (role) {
      filteredList = filteredList.filter((item) => item.role === role);
    }

    // 分页
    const start = (page - 1) * page_size;
    const end = start + page_size;
    const usersData = filteredList.slice(start, end);

    res.json({
      code: 200,
      message: 'success',
      data: {
        users: usersData,
        total: filteredList.length,
      },
    });
  },
  // IAM 创建用户接口
  'POST /v1/iam/users': (req: any, res: any) => {
    const { username, password, email, role, avatar } = req.body;
    const newUser = {
      id: userList.length + 1,
      username,
      email: email || '',
      avatar: avatar || '',
      role: role || 'user',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    userList.push(newUser);

    res.json({
      code: 200,
      message: 'success',
      data: newUser,
    });
  },
  // IAM 更新用户接口
  'PUT /v1/iam/users/:id': (req: any, res: any) => {
    const { id } = req.params;
    const { username, email, role, avatar } = req.body;
    const userIndex = userList.findIndex((u) => u.id === parseInt(id));

    if (userIndex === -1) {
      res.status(404).json({
        code: 404,
        message: '用户不存在',
      });
      return;
    }

    userList[userIndex] = {
      ...userList[userIndex],
      username: username || userList[userIndex].username,
      email: email !== undefined ? email : userList[userIndex].email,
      role: role || userList[userIndex].role,
      avatar: avatar !== undefined ? avatar : userList[userIndex].avatar,
      updated_at: new Date().toISOString(),
    };

    res.json({
      code: 200,
      message: 'success',
      data: userList[userIndex],
    });
  },
  // IAM 删除用户接口
  'DELETE /v1/iam/users/:id': (req: any, res: any) => {
    const { id } = req.params;
    const userIndex = userList.findIndex((u) => u.id === parseInt(id));

    if (userIndex === -1) {
      res.status(404).json({
        code: 404,
        message: '用户不存在',
      });
      return;
    }

    userList.splice(userIndex, 1);

    res.json({
      code: 200,
      message: '删除成功',
    });
  },
  // IAM 修改用户密码接口
  'PUT /v1/iam/users/:id/password': (req: any, res: any) => {
    const { id } = req.params;
    const { newPassword } = req.body;
    const user = userList.find((u) => u.id === parseInt(id));

    if (!user) {
      res.status(404).json({
        code: 404,
        message: '用户不存在',
      });
      return;
    }

    res.json({
      code: 200,
      message: '密码修改成功',
    });
  },
};
